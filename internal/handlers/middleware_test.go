package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_service "github.com/HeadGardener/coursework/internal/handlers/mocks"
	"github.com/HeadGardener/coursework/internal/lib/auth"
	"github.com/HeadGardener/coursework/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestIdentifyUserMiddleware(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthService, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthService, token string) {
				s.EXPECT().ParseAccessToken(token).Return(auth.UserAttributes{
					ID:   "1",
					Role: models.RoleUser,
					Age:  20,
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"1"`,
		},
		{
			name:                 "no header",
			headerName:           "",
			mockBehavior:         func(s *mock_service.MockAuthService, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Msg":"failed while identifying user","Error":"empty auth header"}`,
		},
		{
			name:                 "invalid bearer",
			headerName:           "Authorization",
			headerValue:          "Barer token",
			mockBehavior:         func(s *mock_service.MockAuthService, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Msg":"failed while identifying user","Error":"invalid auth header Barer, must be Bearer"}`,
		},
		{
			name:                 "no token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			mockBehavior:         func(s *mock_service.MockAuthService, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Msg":"failed while identifying user","Error":"jwt token is empty"}`,
		},
		{
			name:        "service failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthService, token string) {
				s.EXPECT().ParseAccessToken(token).Return(auth.UserAttributes{}, errors.New(""))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Msg":"failed while parsing token","Error":""}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authService := mock_service.NewMockAuthService(c)
			tc.mockBehavior(authService, tc.token)

			handler := NewHandler(authService, nil)

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			router.Use(gin.Recovery())
			router.Use(handler.identifyUser)
			router.POST("/protected", gin.HandlerFunc(func(c *gin.Context) {
				c.JSON(http.StatusOK, "1")
			}))

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/protected", nil)
			r.Header.Set(tc.headerName, tc.headerValue)

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
