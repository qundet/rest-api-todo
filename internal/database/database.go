package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dbURL)

	if err != nil {
		return nil, err // TODO: fmt.Errorf add some context
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	
	return db, nil
}