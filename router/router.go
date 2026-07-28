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

	// Routes
	articleGroup := r.Group("/article")
	{
		articleGroup.POST("", articleHandler.CreateArticle)
		articleGroup.POST("/", articleHandler.CreateArticle)

		articleGroup.GET("/:limit/:offset", articleHandler.GetArticles)
		articleGroup.GET("/:id", articleHandler.GetArticleByID)

		// Support PATCH, PUT, and POST for article update as mentioned in test requirements
		articleGroup.PATCH("/:id", articleHandler.UpdateArticle)
		articleGroup.PUT("/:id", articleHandler.UpdateArticle)
		
		// Support DELETE and POST for article deletion
		articleGroup.DELETE("/:id", articleHandler.DeleteArticle)
	}

	return r
}
