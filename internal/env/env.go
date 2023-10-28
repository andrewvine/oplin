package env

import (
	"fmt"
	"log"
	"oplin/internal/projectpath"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Setup() {
	err := godotenv.Load(fmt.Sprintf("%s/.env", projectpath.Root))
	if err != nil {
		log.Println("No .env file found")
	}
}

func getEnv(k, v string) string {
	s := os.Getenv(k)
	if s == "" {
		return v
	}
	return s

}

func BuildDSN(host, user, dbname, password string, port int) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		getEnv("OPLIN_DB_HOST", host),
		getEnv("OPLIN_DB_USER", user),
		getEnv("OPLIN_DB_PASSWORD", password),
		getEnv("OPLIN_DB_NAME", dbname),
		getEnv("OPLIN_DB_PORT", strconv.Itoa(port)))
}

func GetTestDSN() string {
	return ""
}