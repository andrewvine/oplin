package env

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Setup() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
		if os.Getenv("DB_HOST") == "" {
			log.Println("No environment variables like DB_HOST found")
		}
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
		getEnv("DB_HOST", host),
		getEnv("DB_USER", user),
		getEnv("DB_PASSWORD", password),
		getEnv("DB_NAME", dbname),
		getEnv("DB_PORT", strconv.Itoa(port)))
}

func GetTestDSN() string {
	return ""
}

func PrependProjectPath(s string) string {
	path := os.Getenv("OPLIN_PATH")
	if path == "" {
		log.Fatal("OPLIN_PATH not set")
	}
	return fmt.Sprintf("%s/%s", path, s)
}
