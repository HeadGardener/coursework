package service

import (
	"context"

	"github.com/HeadGardener/coursework/internal/models"
)

type DrinkStorage interface {
	GetAll(ctx context.Context, adult bool) ([]models.Drink, error)
	GetByID(ctx context.Context, id int, adult bool) (models.Drink, error)
	Create(ctx context.Context, drink *models.Drink) (int, error)
	Update(ctx context.Context, id int, drink *models.Drink) error
	Delete(ctx context.Context, id int) error
}

type DrinkService struct {
	drinkStorage DrinkStorage
}

func NewDrinkService(drinkStorage DrinkStorage) *DrinkService {
	return &DrinkService{drinkStorage: drinkStorage}
}

func (s *DrinkService) GetAll(ctx context.Context, adult bool) ([]models.Drink, error) {
	return s.drinkStorage.GetAll(ctx, adult)
}

func (s *DrinkService) GetByID(ctx context.Context, id int, adult bool) (models.Drink, error) {
	return s.drinkStorage.GetByID(ctx, id, adult)
}

func (s *DrinkService) Add(ctx context.Context, drink *models.Drink) (int, error) {
	return s.drinkStorage.Create(ctx, drink)
}

func (s *DrinkService) Update(ctx context.Context, id int, drinkInput *models.Drink) error {
	drink, err := s.drinkStorage.GetByID(ctx, id, true)
	if err != nil {
		return err
	}

	if drink.Name != drinkInput.Name {
		drink.Name = drinkInput.Name
	}

	if drink.Type != drinkInput.Type {
		drink.Type = drinkInput.Type
	}

	if drink.Bottle != drinkInput.Bottle && drinkInput.Bottle != 0 {
		drink.Bottle = drinkInput.Bottle
	}

	if drink.Cost != drinkInput.Cost && drinkInput.Cost != 0 {
		drink.Cost = drinkInput.Cost
	}

	return s.drinkStorage.Update(ctx, id, &drink)
}

func (s *DrinkService) Delete(ctx context.Context, id int) error {
	return s.drinkStorage.Delete(ctx, id)
}
