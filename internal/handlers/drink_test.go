package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_service "github.com/HeadGardener/coursework/internal/handlers/mocks"
	"github.com/HeadGardener/coursework/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestAddDrinkHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockDrinkService, drink *models.Drink)

	testTable := []struct {
		name                 string
		inputBody            string
		drink                *models.Drink
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			inputBody: `{
          					"name": "test",
               				"type": "test",
                   			"bottle": 100,
                      		"cost": 100,
                        	"soft": true
                      	}`,
			drink: &models.Drink{
				ID:     0,
				Name:   "test",
				Type:   "test",
				Bottle: 100,
				Cost:   100,
				Soft:   true,
			},
			mockBehavior: func(s *mock_service.MockDrinkService, drink *models.Drink) {
				s.EXPECT().Add(gomock.Any(), drink).Return(0, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":0}`,
		},
		{
			name: "invalid bottle",
			inputBody: `{
          					"name": "test",
               				"type": "test",
                   			"bottle": 0,
                      		"cost": 100,
                        	"soft": true
                      	}`,
			mockBehavior:         func(s *mock_service.MockDrinkService, drink *models.Drink) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"Msg":"failed while validating drink request","Error":"invalid bottle: bottle can't be less or equals 0"}`,
		},
		{
			name: "invalid cost",
			inputBody: `{
          					"name": "test",
               				"type": "test",
                   			"bottle": 100,
                      		"cost": -1,
                        	"soft": true
                      	}`,
			mockBehavior:         func(s *mock_service.MockDrinkService, drink *models.Drink) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"Msg":"failed while validating drink request","Error":"invalid cost: cost can't be less than 0"}`,
		},
		{
			name: "service failure",
			inputBody: `{
          					"name": "test",
               				"type": "test",
                   			"bottle": 100,
                      		"cost": 100,
                        	"soft": true
                      	}`,
			drink: &models.Drink{
				Name:   "test",
				Type:   "test",
				Bottle: 100,
				Cost:   100,
				Soft:   true,
			},
			mockBehavior: func(s *mock_service.MockDrinkService, drink *models.Drink) {
				s.EXPECT().Add(gomock.Any(), drink).Return(0, errors.New(""))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"Msg":"failed while adding drink","Error":""}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			drink := mock_service.NewMockDrinkService(c)
			tc.mockBehavior(drink, tc.drink)

			handler := NewHandler(nil, drink)

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			router.Use(gin.Recovery())
			router.POST("/api/drinks/", handler.addDrink)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/drinks/", bytes.NewBufferString(tc.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
