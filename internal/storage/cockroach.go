package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/HeadGardener/coursework/internal/lib/hash"
	"github.com/HeadGardener/coursework/internal/models"
	"github.com/google/uuid"

	"github.com/HeadGardener/coursework/internal/config"
	_ "github.com/cockroachdb/cockroach-go/v2/crdb"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

//nolint:gomnd
func initTable(ctx context.Context, tx *sql.Tx) error {
	log.Println("drop existing users table if necessary")
	if _, err := tx.ExecContext(ctx, `drop table if exists users`); err != nil {
		return err
	}

	log.Println("drop existing drinks table if necessary")
	if _, err := tx.ExecContext(ctx, `drop table if exists drinks`); err != nil {
		return err
	}

	// create users table
	log.Println("creating users table")
	if _, err := tx.ExecContext(ctx,
		`create table users (
    				id uuid primary key,                                     
            		username varchar(255) not null unique,                   
                    name varchar(255) not null,                              
                    role integer not null,
                    age integer not null,                             
                    password_hash varchar(255) not null
                   );`); err != nil {
		return err
	}

	log.Println("inserting admin into users table")
	if _, err := tx.ExecContext(ctx, `insert into users (id, username, name, role, age, password_hash)
											values ($1, $2, $3, $4, $5, $6)`,
		uuid.NewString(),
		"superadmin",
		"admin",
		models.RoleAdmin,
		30,
		hash.GetPasswordHash("1234567890")); err != nil {
		return err
	}

	// create drinks table
	log.Println("creating drinks table")
	if _, err := tx.ExecContext(ctx,
		`create table drinks (
    				id serial primary key,                                     
            		name varchar(255) not null unique,
            		type varchar(255) not null,
            		bottle integer not null default 1000,
            		cost float not null,
            		is_soft bool not null
                   );`); err != nil {
		return err
	}

	log.Println("inserting drinks into drinks table")
	if _, err := tx.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"VOSS",
		"water",
		700,
		10,
		true); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"Dr.Pepper",
		"soda",
		300,
		3,
		true); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"Mountain Dew",
		"soda",
		500,
		2,
		true); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"Corona Extra",
		"beer",
		355,
		5,
		false); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"Jagermeister",
		"liquor",
		1000,
		40,
		false); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `insert into drinks (name, type, bottle, cost, is_soft)
											values ($1,$2,$3,$4,$5)`,
		"Maker's Mark",
		"bourbon",
		1000,
		50,
		false); err != nil {
		return err
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

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err = initTable(ctx, tx); err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("unexpected error: unable to rollback: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("unexpected error: unable to commit: %w", err)
	}

	return db, nil
}
