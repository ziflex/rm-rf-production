package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/ziflex/dbx"
)

func New(opts Options) (dbx.Database, error) {
	db, err := sql.Open("postgres", toConnectionString(opts))

	if err != nil {
		return nil, err
	}

	return db, nil
}
