package handler

import (
	"net/http"
	"strconv"

	"sharing-vision-backend/dto"
	"sharing-vision-backend/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ArticleHandler struct {
	service service.ArticleService
}

func NewArticleHandler(service service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

// Helper to format validation errors
func formatValidationError(err error) map[string]string {
	errs := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := e.Field()
			switch field {
			case "Title":
				errs["title"] = "Title is required and must be at least 20 characters long."
			case "Content":
				errs["content"] = "Content is required and must be at least 200 characters long."
			case "Category":
				errs["category"] = "Category is required and must be at least 3 characters long."
			case "Status":
				errs["status"] = "Status is required and must be one of: publish, draft, thrash."
			default:
				errs[field] = e.Error()
			}
		}
	} else {
		errs["error"] = err.Error()
	}
	return errs
}

// POST /article/
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validation failed",
			"errors":  formatValidationError(err),
		})
		return
	}

	res, err := h.service.CreateArticle(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// GET /article/:limit/:offset
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	limitStr := c.Param("limit")
	offsetStr := c.Param("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	res, err := h.service.GetArticles(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Returning list directly as specified in PDF specification array [{ "title": "", ... }] or data object
	c.JSON(http.StatusOK, res.Data)
}

// GET /article/:id
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid article ID"})
		return
	}

	res, err := h.service.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// PATCH/PUT/POST /article/:id
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid article ID"})
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validation failed",
			"errors":  formatValidationError(err),
		})
		return
	}

	res, err := h.service.UpdateArticle(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DELETE/POST /article/:id
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid article ID"})
		return
	}

	if err := h.service.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Article status updated to thrash successfully",
	})
}
