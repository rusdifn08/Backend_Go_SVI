package main

import (
	"log"
	"os"

	"sharing-vision-backend/config"
	"sharing-vision-backend/handler"
	"sharing-vision-backend/repository"
	"sharing-vision-backend/router"
	"sharing-vision-backend/service"
)

func main() {
	log.Println("Starting Sharing Vision Article Microservice...")

	// Initialize Database (GORM + MySQL)
	db := config.InitDB()

	// Initialize Layers (Clean Architecture)
	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	// Setup Router
	r := router.SetupRouter(articleHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
