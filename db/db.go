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

	template := "%s:%s@tcp(%s:3306)/%s?charset=utf8mb4"

	dsn := fmt.Sprintf(
		template,
		"root",
		utils.Getenv("password", "atr"),
		utils.Getenv("host", "localhost"),
		utils.Getenv("dbname", "entretien"),
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
