package router

import (
	"sharing-vision-backend/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(articleHandler *handler.ArticleHandler) *gin.Engine {
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Routes for /article
	articleGroup := r.Group("/article")
	{
		// Create article
		articleGroup.POST("", articleHandler.CreateArticle)
		articleGroup.POST("/", articleHandler.CreateArticle)

		// Pagination GET /article/:id/:offset (where :id represents limit)
		// Sharing the wildcard parameter name ':id' avoids Gin wildcard routing conflict
		articleGroup.GET("/:id/:offset", articleHandler.GetArticles)

		// Get detail GET /article/:id
		articleGroup.GET("/:id", articleHandler.GetArticleByID)

		// Support PATCH, PUT for article update
		articleGroup.PATCH("/:id", articleHandler.UpdateArticle)
		articleGroup.PUT("/:id", articleHandler.UpdateArticle)

		// Support DELETE for article deletion (soft delete to thrash)
		articleGroup.DELETE("/:id", articleHandler.DeleteArticle)
	}

	return r
}
