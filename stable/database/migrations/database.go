package migrations

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"../entities"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
