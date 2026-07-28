package service

import (
	"strings"

	"sharing-vision-backend/domain"
	"sharing-vision-backend/dto"
	"sharing-vision-backend/repository"
)

type ArticleService interface {
	CreateArticle(req dto.CreateArticleRequest) (*dto.ArticleResponse, error)
	GetArticles(limit, offset int) (*dto.PaginationResponse, error)
	GetArticleByID(id int) (*dto.ArticleResponse, error)
	UpdateArticle(id int, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error)
	DeleteArticle(id int) error
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func normalizeStatus(status string) domain.ArticleStatus {
	s := strings.ToLower(strings.TrimSpace(status))
	switch s {
	case "publish":
		return domain.StatusPublish
	case "draft":
		return domain.StatusDraft
	case "thrash", "trash":
		return domain.StatusThrash
	default:
		return domain.StatusDraft
	}
}

func (s *articleService) CreateArticle(req dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	article := &domain.Article{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   normalizeStatus(req.Status),
	}

	if err := s.repo.Create(article); err != nil {
		return nil, err
	}

	return s.toArticleResponse(article), nil
}

func (s *articleService) GetArticles(limit, offset int) (*dto.PaginationResponse, error) {
	articles, total, err := s.repo.GetPaginated(limit, offset)
	if err != nil {
		return nil, err
	}

	var responseList []dto.ArticleResponse
	for _, art := range articles {
		responseList = append(responseList, *s.toArticleResponse(&art))
	}

	if responseList == nil {
		responseList = []dto.ArticleResponse{}
	}

	return &dto.PaginationResponse{
		Data:   responseList,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *articleService) GetArticleByID(id int) (*dto.ArticleResponse, error) {
	article, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toArticleResponse(article), nil
}

func (s *articleService) UpdateArticle(id int, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	existing.Title = req.Title
	existing.Content = req.Content
	existing.Category = req.Category
	existing.Status = normalizeStatus(req.Status)

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return s.toArticleResponse(existing), nil
}

func (s *articleService) DeleteArticle(id int) error {
	return s.repo.Delete(id)
}

func (s *articleService) toArticleResponse(article *domain.Article) *dto.ArticleResponse {
	return &dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		Content:     article.Content,
		Category:    article.Category,
		CreatedDate: article.CreatedDate.Format("2006-01-02 15:04:05"),
		UpdatedDate: article.UpdatedDate.Format("2006-01-02 15:04:05"),
		Status:      string(article.Status),
	}
}
