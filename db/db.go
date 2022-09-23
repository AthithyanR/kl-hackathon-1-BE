package db

import (
	"fmt"
	"log"

	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDb() {

	template := "%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true"

	dsn := fmt.Sprintf(
		template,
		utils.Getenv("DB_USER", "root"),
		utils.Getenv("DB_PASS", "password"),
		utils.Getenv("DB_HOST", "localhost"),
		utils.Getenv("DB_NAME", "db"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(`Unable to establish database connection :-(`)
	}

	log.Println("Database connection established")
	DB = db
}
