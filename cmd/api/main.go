package main

import (
	"github.com/joho/godotenv"
	"stable/database/migrations"
	"log"
	"os"
)
	
func main() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}

	migrations.ConnectDB()
}