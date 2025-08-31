package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ziflex/dbx"
)

func New(opts Options) (dbx.Database, error) {
	db, err := sqlx.Open("postgres", toConnectionString(opts))

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
