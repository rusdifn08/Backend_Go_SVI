package repository

import (
	"errors"
	"strings"
	"sync"
	"time"

	"sharing-vision-backend/domain"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(article *domain.Article) error
	GetPaginated(limit, offset int) ([]domain.Article, int64, error)
	GetByID(id int) (*domain.Article, error)
	Update(article *domain.Article) error
	Delete(id int) error
}

type articleRepository struct {
	db *gorm.DB
	// In-memory fallback if database connection is not established during initial testing
	memStore map[int]*domain.Article
	mu       sync.RWMutex
	nextID   int
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	repo := &articleRepository{
		db:       db,
		memStore: make(map[int]*domain.Article),
		nextID:   1,
	}

	// Pre-seed mock data if using in-memory fallback
	if db == nil {
		repo.seedMockData()
	}

	return repo
}

func (r *articleRepository) seedMockData() {
	now := time.Now()
	mockArticles := []domain.Article{
		{
			ID:          1,
			Title:       "Panduan Lengkap Microservice dengan Golang dan Clean Architecture",
			Content:     "Microservice architecture telah menjadi standar industri untuk membangun aplikasi scalable. Dalam panduan ini, kita akan membahas penerapan Clean Architecture menggunakan bahasa pemrograman Go (Golang), framework Gin, dan GORM untuk manajemen database MySQL secara efisien.",
			Category:    "Technology",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusPublish,
		},
		{
			ID:          2,
			Title:       "Mengenal Next.js App Router dan State Management Zustand",
			Content:     "Next.js App Router membawa perubahan mendasar dalam cara kita membangun frontend modern. Dipadukan dengan Zustand sebagai state manager yang ringan dan intuitif, pengembang dapat mengelola global state tanpa overhead kompleksitas Redux.",
			Category:    "Frontend",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusDraft,
		},
		{
			ID:          3,
			Title:       "Artikel Contoh Yang Sudah Dihapus Ke Dalam Status Thrash",
			Content:     "Ini adalah konten artikel contoh yang dipindahkan ke status thrash untuk pengujian fitur restore atau permanent removal pada dashboard manajemen post.",
			Category:    "General",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusThrash,
		},
	}

	for _, a := range mockArticles {
		art := a
		r.memStore[art.ID] = &art
		if art.ID >= r.nextID {
			r.nextID = art.ID + 1
		}
	}
}

func (r *articleRepository) Create(article *domain.Article) error {
	if r.db != nil {
		article.CreatedDate = time.Now()
		article.UpdatedDate = time.Now()
		return r.db.Create(article).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	article.ID = r.nextID
	r.nextID++
	article.CreatedDate = time.Now()
	article.UpdatedDate = time.Now()
	r.memStore[article.ID] = article
	return nil
}

func (r *articleRepository) GetPaginated(limit, offset int) ([]domain.Article, int64, error) {
	if r.db != nil {
		var articles []domain.Article
		var total int64

		if err := r.db.Model(&domain.Article{}).Count(&total).Error; err != nil {
			return nil, 0, err
		}

		query := r.db.Order("id DESC")
		if limit > 0 {
			query = query.Limit(limit)
		}
		if offset >= 0 {
			query = query.Offset(offset)
		}

		if err := query.Find(&articles).Error; err != nil {
			return nil, 0, err
		}
		return articles, total, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []domain.Article
	for _, a := range r.memStore {
		all = append(all, *a)
	}

	total := int64(len(all))
	if offset >= len(all) {
		return []domain.Article{}, total, nil
	}

	end := offset + limit
	if limit <= 0 || end > len(all) {
		end = len(all)
	}

	return all[offset:end], total, nil
}

func (r *articleRepository) GetByID(id int) (*domain.Article, error) {
	if r.db != nil {
		var article domain.Article
		if err := r.db.First(&article, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("article not found")
			}
			return nil, err
		}
		return &article, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	article, ok := r.memStore[id]
	if !ok {
		return nil, errors.New("article not found")
	}
	return article, nil
}

func (r *articleRepository) Update(article *domain.Article) error {
	if r.db != nil {
		article.UpdatedDate = time.Now()
		return r.db.Save(article).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.memStore[article.ID]
	if !ok {
		return errors.New("article not found")
	}

	article.CreatedDate = existing.CreatedDate
	article.UpdatedDate = time.Now()
	r.memStore[article.ID] = article
	return nil
}

func (r *articleRepository) Delete(id int) error {
	if r.db != nil {
		// Soft delete: update status to "thrash" as specified in PDF/reqs or perform update
		return r.db.Model(&domain.Article{}).Where("id = ?", id).Update("status", strings.ToLower(string(domain.StatusThrash))).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	article, ok := r.memStore[id]
	if !ok {
		return errors.New("article not found")
	}

	article.Status = domain.StatusThrash
	article.UpdatedDate = time.Now()
	return nil
}
