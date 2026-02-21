package migrations

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"fmt"
	"os"
	"stable/database/entities"
)

var DB *gorm.DB

func ConnectDB() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")

	log.Println("Running auto-migrations...")


	log.Println("🚀 Running auto-migrations...")
	
	err = DB.AutoMigrate(
    &entities.User{},
    &entities.ProfileUser{},
    &entities.Exercise{},
    &entities.ExerciseContent{},
    &entities.DailyChallenge{},
    &entities.DailyChallengeExercise{},
    &entities.UserDailyActivityLog{},
    &entities.UserStreak{},
	)

	if err != nil {
		log.Fatalf("Failed to run auto-migrations: %v", err)
	}
	log.Println("Auto-migrations completed successfully!")
}

func GetDB() *gorm.DB {
	return DB
}