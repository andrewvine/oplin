package utils

import (
	"oplin/internal/env"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func GetTestDB() *sql.DB {
	dsn := env.GetTestDSN()
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not get test db[%s]", err)
	}
	return DB
}
