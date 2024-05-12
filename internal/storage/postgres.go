package storage

import (
	"context"
	"log"

	"github.com/HeadGardener/coursework/internal/config"
	"github.com/HeadGardener/coursework/internal/lib/hash"
	"github.com/HeadGardener/coursework/internal/models"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

//nolint:gomnd
func initTable(ctx context.Context, db *sqlx.DB) error {
	log.Println("inserting admin into users table")
	if _, err := db.ExecContext(ctx, `insert into users (id, username, name, role, age, password_hash)
											values ($1, $2, $3, $4, $5, $6)`,
		uuid.NewString(),
		"superadmin",
		"admin",
		models.RoleAdmin,
		30,
		hash.GetStringHash("1234567890")); err != nil {
		log.Println("failed to insert admin: ", err.Error())
	}

	log.Println("inserting drinks into drinks table")
	if _, err := db.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"VOSS",
		"water",
		700,
		10,
		true); err != nil {
		log.Println("failed to insert drink while initializing: ", err.Error())
	}

	if _, err := db.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"Dr.Pepper",
		"soda",
		300,
		3,
		true); err != nil {
		log.Println("failed to insert drink while initializing: ", err.Error())
	}

	if _, err := db.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"Mountain Dew",
		"soda",
		500,
		2,
		true); err != nil {
		log.Println("failed to insert drink while initializing: ", err.Error())
	}

	if _, err := db.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"Corona Extra",
		"beer",
		355,
		5,
		false); err != nil {
		log.Println("failed to insert drink while initializing: ", err.Error())
	}

	if _, err := db.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"Jagermeister",
		"liquor",
		1000,
		40,
		false); err != nil {
		log.Println("failed to insert drink while initializing: ", err.Error())
	}

	if _, err := db.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"Maker's Mark",
		"bourbon",
		1000,
		50,
		false); err != nil {
		log.Println("failed to insert drink while initializing: ", err.Error())
	}

	return nil
}

func NewDB(ctx context.Context, conf config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", conf.URL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	initTable(ctx, db)

	return db, nil
}
