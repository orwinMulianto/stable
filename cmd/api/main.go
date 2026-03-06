package main

import (
	"log"
	"os"
	"stable/database/migrations"
	"stable/modules/auth"
	"stable/modules/user"
	"stable/packages/utils"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
	
func main() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("No .env file found")
		}
	}

	migrations.ConnectDB()
	utils.InitLogger()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	v1 := api.Group("/v1")

	{
		auth.AuthRouter(v1)
		user.UserRouter(v1)
	}

	err := router.Run(":8080")
	if err != nil {
    log.Fatal(err)
	}

}