package config

import (
	"fmt"
	"log"
	"os"

	"sharing-vision-backend/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "article")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("Warning: Failed to connect to MySQL database: %v. Running in mock/fallback mode or check DB connection.", err)
		return nil
	}

	log.Println("Database connection established successfully.")

	// Auto Migrate schema
	if err := db.AutoMigrate(&domain.Article{}); err != nil {
		log.Printf("AutoMigrate warning: %v", err)
	}

	DB = db
	return db
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
