package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func getEnv(k, v string) string {
	s := os.Getenv(k)
	if s == "" {
		return v
	}
	return s

}

func buildDSN(host, user, dbname, password string, port int, sslmode string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("OPLIN_TEST_DB_HOST", host),
		getEnv("OPLIN_TEST_DB_USER", user),
		getEnv("OPLIN_TEST_DB_PASSWORD", password),
		getEnv("OPLIN_TEST_DB_NAME", dbname),
		getEnv("OPLIN_TEST_DB_PORT", strconv.Itoa(port)),
		getEnv("OPLIN_TEST_DB_SSLMODE", sslmode))
}

func GetTestDB() *sql.DB {
	dsn := buildDSN("localhost", "test_oplin", "test_oplin", "topsecret", 5432, "disable")
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not get test db[%s]", err)
	}
	return DB
}
