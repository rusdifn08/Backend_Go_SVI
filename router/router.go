package router

import (
	"sharing-vision-backend/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(articleHandler *handler.ArticleHandler) *gin.Engine {
	r := gin.Default()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Root Health Check Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Sharing Vision Article Microservice API is running smoothly",
		})
	})

	// Routes for /article
	articleGroup := r.Group("/article")
	{
		articleGroup.POST("", articleHandler.CreateArticle)
		articleGroup.POST("/", articleHandler.CreateArticle)

		articleGroup.GET("/:id/:offset", articleHandler.GetArticles)
		articleGroup.GET("/:id", articleHandler.GetArticleByID)

		articleGroup.PATCH("/:id", articleHandler.UpdateArticle)
		articleGroup.PUT("/:id", articleHandler.UpdateArticle)

		articleGroup.DELETE("/:id", articleHandler.DeleteArticle)
	}

	return r
}
