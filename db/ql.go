package db

import (
	"database/sql"

	_ "github.com/cznic/ql/driver"
)

type DB struct {
	*sql.DB
}

func Open(path string) (*DB, error) {
	db, err := sql.Open("ql", path)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB: db,
	}, nil
}
