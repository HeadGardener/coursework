package storage

import (
	"context"

	"github.com/HeadGardener/coursework/internal/models"
	"github.com/jmoiron/sqlx"
)

type DrinkStorage struct {
	db *sqlx.DB
}

func NewDrinkStorage(db *sqlx.DB) *DrinkStorage {
	return &DrinkStorage{db: db}
}

func (s *DrinkStorage) GetAll(ctx context.Context, adult bool) ([]models.Drink, error) {
	var drinks []models.Drink

	var query = `select * from drinks`
	if !adult {
		query += ` where is_soft=true`
	}

	if err := s.db.SelectContext(ctx, &drinks, query); err != nil {
		return nil, err
	}

	return drinks, nil
}

func (s *DrinkStorage) GetByID(ctx context.Context, id int, adult bool) (models.Drink, error) {
	var drink models.Drink

	var query = `select * from drinks where id=$1`
	if !adult {
		query += ` and is_soft=true`
	}

	if err := s.db.GetContext(ctx, &drink, query, id); err != nil {
		return models.Drink{}, err
	}

	return drink, nil
}

func (s *DrinkStorage) Create(ctx context.Context, drink *models.Drink) (int, error) {
	var id int

	if err := s.db.QueryRowContext(ctx,
		`insert into drinks (name, type, bottle, cost, is_soft) values($1,$2,$3,$4,$5) returning id`,
		drink.Name, drink.Type, drink.Bottle, drink.Cost, drink.Soft).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *DrinkStorage) Update(ctx context.Context, id int, drink *models.Drink) error {
	if _, err := s.db.ExecContext(ctx, `update drinks set name=$1, type=$2, bottle=$3, cost=$4 where id=$5`,
		drink.Name, drink.Type, drink.Bottle, drink.Cost, id); err != nil {
		return err
	}

	return nil
}

func (s *DrinkStorage) Delete(ctx context.Context, id int) error {
	if _, err := s.db.ExecContext(ctx, `delete from drinks where id=$1`,
		id); err != nil {
		return err
	}

	return nil
}
