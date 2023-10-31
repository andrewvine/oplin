package env

import (
	"fmt"
	"log"
	"oplin/internal/projectpath"

	"github.com/joho/godotenv"
)

func Setup() {
	err := godotenv.Load(fmt.Sprintf("%s/.env", projectpath.Root))
	if err != nil {
		log.Println("No .env file found")
	}
}