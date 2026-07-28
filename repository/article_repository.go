package repository

import (
	"errors"
	"log"
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
	db       *gorm.DB
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

	if db != nil {
		repo.seedTiDBDatabase()
	} else {
		repo.seedMockMemory()
	}

	return repo
}

func (r *articleRepository) getDummyArticles() []domain.Article {
	now := time.Now()
	return []domain.Article{
		{
			Title:       "Panduan Lengkap Microservice dengan Golang dan Clean Architecture",
			Content:     "Microservice architecture telah menjadi standar industri untuk membangun aplikasi scalable. Dalam panduan ini, kita akan membahas penerapan Clean Architecture menggunakan bahasa pemrograman Go (Golang), framework Gin, dan GORM untuk manajemen database TiDB Cloud secara efisien dan performan tinggi.",
			Category:    "Technology",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusPublish,
		},
		{
			Title:       "Mengenal Next.js App Router dan State Management Zustand",
			Content:     "Next.js App Router membawa perubahan mendasar dalam cara kita membangun frontend modern. Dipadukan dengan Zustand sebagai state manager yang ringan dan intuitif, pengembang dapat mengelola global state tanpa overhead kompleksitas Redux pada dashboard CMS artikel.",
			Category:    "Frontend",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusPublish,
		},
		{
			Title:       "Optimasi Query Database MySQL dan TiDB Cloud untuk Aplikasi Scalable",
			Content:     "TiDB Cloud adalah database NewSQL terdistribusi yang kompatibel penuh dengan protokol MySQL. Dalam panduan ini, kita membahas teknik indexing, pemanfaatan connection pooling, dan struktur schema SQL untuk menjaga latency query tetap stabil dalam hitungan milidetik.",
			Category:    "Database",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusPublish,
		},
		{
			Title:       "Rancangan Arsitektur Sistem E-Commerce dengan Microservices Golang",
			Content:     "Membangun sistem e-commerce berskala besar membutuhkan pemisahan concern pada level service. Artikel draft ini membahas strategi pemisahan service katalog, payment gateway, manajemen inventoris, dan autentikasi JWT di dalam ekosistem microservices Go.",
			Category:    "Architecture",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusDraft,
		},
		{
			Title:       "Panduan Lengkap Integrasi CI/CD Pipeline dengan Docker dan Kubernetes",
			Content:     "Otomatisasi pengujian dan penggelaran aplikasi adalah kunci kecepatan inovasi produk. Dalam artikel draft ini, kita mendiskusikan langkah konfigurasi GitHub Actions, pembentukan image Docker yang teroptimasi, dan penggelaran otomatis ke kluster Kubernetes.",
			Category:    "DevOps",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusDraft,
		},
		{
			Title:       "Draft Artikel Mengenai Best Practices Security pada RESTful API",
			Content:     "Keamanan API merupakan aspek kritis yang tidak boleh diabaikan. Penulisan draf artikel ini mengulas implementasi Rate Limiting, Input Sanitization untuk mencegah SQL Injection, CORS Policy yang ketat, serta enkripsi payload sensitif.",
			Category:    "Security",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusDraft,
		},
		{
			Title:       "Artikel Lama Mengenai Monolithic Framework Yang Sudah Dihapus",
			Content:     "Ini adalah konten artikel contoh yang dipindahkan ke status thrash untuk pengujian fitur restore atau permanent removal pada dashboard manajemen post Sharing Vision.",
			Category:    "General",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusThrash,
		},
		{
			Title:       "Catatan Riset Perbandingan Performance Redis vs Memcached",
			Content:     "Arsip artikel perbandingan caching layer yang sudah dipindahkan ke kategori sampah (thrash) setelah dilakukan revisi strategi in-memory storage pada infrastruktur backend.",
			Category:    "Backend",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusThrash,
		},
	}
}

func (r *articleRepository) seedTiDBDatabase() {
	var count int64
	if err := r.db.Model(&domain.Article{}).Count(&count).Error; err == nil {
		if count == 0 {
			log.Println("Seeding initial dummy data to TiDB Cloud database...")
			dummies := r.getDummyArticles()
			for _, article := range dummies {
				if err := r.db.Create(&article).Error; err != nil {
					log.Printf("Error seeding article: %v", err)
				}
			}
			log.Println("Database TiDB Cloud successfully seeded with dummy data!")
		}
	}
}

func (r *articleRepository) seedMockMemory() {
	dummies := r.getDummyArticles()
	for i, article := range dummies {
		art := article
		art.ID = i + 1
		r.memStore[art.ID] = &art
		r.nextID = art.ID + 1
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
