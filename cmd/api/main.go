package main

import (
	"log"
	"os"
	"stable/database/migrations"
	"stable/modules/auth"
	"stable/modules/users"
	"stable/packages/utils"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))
	
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("No .env file found")
		}
	}
migrations.ConnectDB()
utils.InitLogger()
router.SetTrustedProxies(nil)

api := router.Group("/api")
	v1 := api.Group("/v1")

	{
		auth.AuthRouter(v1)
		users.UserRouter(v1)
	}

	err := router.Run(":8080")
	if err != nil {
    log.Fatal(err)
	}

}