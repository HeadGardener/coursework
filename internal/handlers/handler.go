package handlers

import (
	"context"
	"net/http"

	"github.com/HeadGardener/coursework/internal/lib/auth"
	"github.com/HeadGardener/coursework/internal/models"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	SignUp(ctx context.Context, username, name string, age int, password string) (string, error)
	SignIn(ctx context.Context, username, password string) (models.Tokens, error)
	ParseAccessToken(token string) (auth.UserAttributes, error)
	Refresh(ctx context.Context, accessToken, refreshToken string) (models.Tokens, error)
	LogOut(ctx context.Context, userID string) error
}

type DrinkService interface {
	GetAll(ctx context.Context, adult bool) ([]models.Drink, error)
	GetByID(ctx context.Context, id int, adult bool) (models.Drink, error)
	Add(ctx context.Context, drink *models.Drink) (int, error)
	Update(ctx context.Context, id int, drink *models.Drink) error
	Delete(ctx context.Context, id int) error
}

type Handler struct {
	authService  AuthService
	drinkService DrinkService
}

func NewHandler(authService AuthService, drinkService DrinkService) *Handler {
	return &Handler{
		authService:  authService,
		drinkService: drinkService,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	router := gin.New()

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)
			auth.POST("/refresh", h.refresh)
			auth.PUT("/logout", h.identifyUser, h.logout)
		}

		drinks := api.Group("/drinks", h.identifyUser, h.checkAge)
		{
			drinks.GET("/", h.viewDrinks)
			drinks.GET("/:id", h.viewByID)
			drinks.POST("/", h.identifyRole, h.addDrink)
			drinks.PUT("/:id", h.identifyRole, h.updateDrink)
			drinks.DELETE("/:id", h.identifyRole, h.deleteDrink)
		}
	}

	return router
}
