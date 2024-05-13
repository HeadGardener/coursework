package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HeadGardener/coursework/internal/dto"
	mock_service "github.com/HeadGardener/coursework/internal/handlers/mocks"
	"github.com/HeadGardener/coursework/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestSignUpHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthService, user dto.SignUpReq)

	testTable := []struct {
		name                 string
		inputBody            string
		user                 dto.SignUpReq
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			inputBody: `{
    						"username": "user",
          					"name": "test",
               				"age": 20,
                   			"password": "testPass"
                      	}`,
			user: dto.SignUpReq{
				Username: "user",
				Name:     "test",
				Age:      20,
				Password: "testPass",
			},
			mockBehavior: func(s *mock_service.MockAuthService, user dto.SignUpReq) {
				s.EXPECT().SignUp(gomock.Any(), user.Username, user.Name, user.Age, user.Password).Return("1", nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":"1"}`,
		},
		{
			name: "invalid name",
			inputBody: `{
    						"username": "user",
          					"name": "tes3t",
               				"age": 20,
                   			"password": "testPass"
                      	}`,
			mockBehavior:         func(s *mock_service.MockAuthService, user dto.SignUpReq) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"Msg":"failed while validating sign up request","Error":"invalid name: must contain only letters"}`,
		},
		{
			name: "invalid age",
			inputBody: `{
    						"username": "user",
          					"name": "test",
               				"age": -1,
                   			"password": "testPass"
                      	}`,
			mockBehavior:         func(s *mock_service.MockAuthService, user dto.SignUpReq) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"Msg":"failed while validating sign up request","Error":"invalid age: can't be less than 0 or greater than 111"}`,
		},
		{
			name: "service failure",
			inputBody: `{
    						"username": "user",
          					"name": "test",
               				"age": 20,
                   			"password": "testPass"
                      	}`,
			user: dto.SignUpReq{
				Username: "user",
				Name:     "test",
				Age:      20,
				Password: "testPass",
			},
			mockBehavior: func(s *mock_service.MockAuthService, user dto.SignUpReq) {
				s.EXPECT().SignUp(gomock.Any(), user.Username, user.Name, user.Age, user.Password).Return("", errors.New(""))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"Msg":"failed while signing up","Error":""}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthService(c)
			tc.mockBehavior(auth, tc.user)

			handler := NewHandler(auth, nil)

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			router.Use(gin.Recovery())
			router.POST("/api/auth/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/auth/sign-up", bytes.NewBufferString(tc.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestSignInHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthService, user dto.SignInReq)

	testTable := []struct {
		name                 string
		inputBody            string
		user                 dto.SignInReq
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			inputBody: `{
    						"username": "user",
                   			"password": "testPass"
                      	}`,
			user: dto.SignInReq{
				Username: "user",
				Password: "testPass",
			},
			mockBehavior: func(s *mock_service.MockAuthService, user dto.SignInReq) {
				s.EXPECT().SignIn(gomock.Any(), user.Username, user.Password).Return(models.Tokens{
					AccessToken:  "token1",
					RefreshToken: "token2",
				}, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"access_token":"token1","refresh_token":"token2"}`,
		},
		{
			name: "invalid username",
			inputBody: `{
    						"username": "user!",
                   			"password": "testPass"
                      	}`,
			mockBehavior:         func(s *mock_service.MockAuthService, user dto.SignInReq) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"Msg":"failed while validating sign in request","Error":"invalid username: must contain only letters, numbers and symbols(_-) "}`,
		},
		{
			name: "service failure",
			inputBody: `{
    						"username": "user",
                   			"password": "testPass"
                      	}`,
			user: dto.SignInReq{
				Username: "user",
				Password: "testPass",
			},
			mockBehavior: func(s *mock_service.MockAuthService, user dto.SignInReq) {
				s.EXPECT().SignIn(gomock.Any(), user.Username, user.Password).Return(models.Tokens{}, errors.New(""))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"Msg":"failed while signing in","Error":""}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthService(c)
			tc.mockBehavior(auth, tc.user)

			handler := NewHandler(auth, nil)

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			router.Use(gin.Recovery())
			router.POST("/api/auth/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/auth/sign-in", bytes.NewBufferString(tc.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
