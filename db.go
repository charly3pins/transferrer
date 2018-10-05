package transferrer

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// NewDB returns a connection to a DB
func NewDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Println("Error connecting to Postgres", err)
		return nil, err
	}

	return db, nil
}
