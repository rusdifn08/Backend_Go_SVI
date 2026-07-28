package dto

type CreateArticleRequest struct {
	Title    string `json:"title" binding:"required,min=20"`
	Content  string `json:"content" binding:"required,min=200"`
	Category string `json:"category" binding:"required,min=3"`
	Status   string `json:"status" binding:"required,oneof=publish draft thrash Publish Draft Thrash"`
}

type UpdateArticleRequest struct {
	Title    string `json:"title" binding:"required,min=20"`
	Content  string `json:"content" binding:"required,min=200"`
	Category string `json:"category" binding:"required,min=3"`
	Status   string `json:"status" binding:"required,oneof=publish draft thrash Publish Draft Thrash"`
}

type ArticleResponse struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	CreatedDate string `json:"created_date,omitempty"`
	UpdatedDate string `json:"updated_date,omitempty"`
	Status      string `json:"status"`
}

type PaginationResponse struct {
	Data  []ArticleResponse `json:"data"`
	Total int64             `json:"total"`
	Limit int               `json:"limit"`
	Offset int              `json:"offset"`
}
