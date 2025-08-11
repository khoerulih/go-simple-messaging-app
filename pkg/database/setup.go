package database

import (
	"fmt"
	"log"

	"github.com/khoerulih/go-simple-messaging-app/app/models"
	"github.com/khoerulih/go-simple-messaging-app/pkg/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.GetEnv("DB_USER", ""),
		env.GetEnv("DB_PASSWORD", ""),
		env.GetEnv("DB_HOST", ""),
		env.GetEnv("DB_PORT", "3307"),
		env.GetEnv("DB_NAME", ""),
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.Logger = logger.Default.LogMode(logger.Info)

	if err := DB.AutoMigrate(&models.User{}, &models.UserSession{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database connection established successfully")

}
