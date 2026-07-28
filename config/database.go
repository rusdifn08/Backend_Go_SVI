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
	dbUser := getEnv("DB_USER", "4TYLPkigyGPqtu5.root")
	dbPass := getEnv("DB_PASSWORD", "eU91Td4D7tdl9YEA")
	dbHost := getEnv("DB_HOST", "gateway01.ap-southeast-1.prod.aws.tidbcloud.com")
	dbPort := getEnv("DB_PORT", "4000")
	dbName := getEnv("DB_NAME", "test")

	// DSN with TLS support for TiDB Cloud
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	log.Printf("Connecting to TiDB Cloud database at %s:%s/%s...", dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("Trying connection without explicit tls param...")
		fallbackDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbPass, dbHost, dbPort, dbName)
		db, err = gorm.Open(mysql.Open(fallbackDsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Printf("Error connecting to TiDB Cloud database: %v", err)
			return nil
		}
	}

	log.Println("Successfully connected to TiDB Cloud database!")

	// Auto Migrate schema for posts table
	if err := db.AutoMigrate(&domain.Article{}); err != nil {
		log.Printf("AutoMigrate error: %v", err)
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
