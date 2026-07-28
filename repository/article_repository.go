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
			Title:       "Pengalaman Mengembangkan Microservice Berbasis Golang dan TiDB Cloud di Environment Production",
			Content:     "Dalam beberapa tahun terakhir, migrasi dari arsitektur monolitik ke microservice menjadi tren utama di banyak perusahaan teknologi. Penggunaan Golang dengan framework Gin terbukti memberikan efisiensi memori yang luar biasa dan latensi yang sangat rendah. Ditambah dengan integrasi TiDB Cloud sebagai NewSQL terdistribusi, kita mendapatkan skalabilitas horizontal tanpa perlu memikirkan kompleksitas sharding database secara manual.",
			Category:    "Teknologi",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusPublish,
		},
		{
			Title:       "Panduan Lengkap Implementasi State Management Zustand Pada Next.js App Router",
			Content:     "Mengelola global state pada Next.js versi 14 dengan App Router terkadang membingungkan bagi pengembang frontend pemula. Zustand hadir sebagai solusi state management yang sangat ringan, tanpa butuh boilerplate yang rumit seperti Redux. Dalam artikel ini, kita akan mengulas bagaimana cara mendefinisikan store Zustand, menangani async API call, serta mengoptimalkan re-render komponen secara efektif.",
			Category:    "Pemrograman",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusPublish,
		},
		{
			Title:       "Prinsip Utama Desain UI/UX Modern: Membangun Antarmuka Kaca Translusen yang Ergonomis",
			Content:     "Tampilan antarmuka yang bersih dan intuitif sangat penting dalam memberikan pengalaman pengguna yang berkesan. Pendekatan desain modern memanfaatkan kontras warna yang tepat, tipografi yang jelas, serta hirarki visual yang terstruktur. Penggunaan aksen warna biru korporat dan layout yang responsif membuat pengoperasian aplikasi CMS menjadi lebih menyenangkan.",
			Category:    "Desain UI/UX",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusPublish,
		},
		{
			Title:       "Strategi Manajemen Database Terdistribusi Menggunakan TiDB Cloud untuk Skalabilitas Tinggi",
			Content:     "Skalabilitas database relational merupakan salah satu tantangan terbesar saat aplikasi berkembang pesat. TiDB Cloud menghadirkan arsitektur NewSQL terdistribusi yang kompatibel dengan protokol MySQL. Fitur auto-scaling dan High Availability bawaan memastikan query tetap dapat berjalan cepat dalam hitungan milidetik meskipun volume data bertambah drastis.",
			Category:    "Database",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusPublish,
		},
		{
			Title:       "Tips Karir Software Engineer: Langkah Menjadi Full-Stack Developer Profesional",
			Content:     "Menjadi seorang Full-Stack Software Engineer yang kompeten membutuhkan pemahaman mendalam tentang konsep dasar software engineering, mulai dari Clean Architecture di sisi backend hingga Atomic Design Pattern di sisi frontend. Selain kemampuan teknis, pemahaman tentang komunikasi dan pemecahan masalah secara terstruktur adalah kunci sukses karir di industri IT.",
			Category:    "Karir IT",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusPublish,
		},
		{
			Title:       "Draf Artikel Mengenai Best Practices Keamanan RESTful API dan Enkripsi Payload",
			Content:     "Keamanan merupakan aspek krusial dalam pembangunan RESTful API modern. Draf artikel ini membahas implementasi Rate Limiting, sanitasi input untuk mencegah XSS dan SQL Injection, pengoperasian protokol HTTPS/TLS, serta autentikasi berbasis JSON Web Token (JWT) secara ketat.",
			Category:    "Keamanan",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusDraft,
		},
		{
			Title:       "Draf Panduan Integrasi Automated Testing dan CI/CD Pipeline Menggunakan GitHub Actions",
			Content:     "Proses pengujian otomatis sebelum deployment membantu menjaga kualitas kode dan mencegah regresi bug di lingkungan production. Draf tutorial ini menjelaskan langkah penyusunan workflow GitHub Actions, eksekusi unit test Golang, dan pembuatan Docker image secara otomatis saat push ke branch main.",
			Category:    "DevOps",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusDraft,
		},
		{
			Title:       "Catatan Riset Internal Mengenai Struktur Data B-Tree dan LSM-Tree Pada Engine Database",
			Content:     "Dokumen riset internal ini menganalisis perbedaan performa antara struktur data B-Tree yang banyak digunakan di database relasional tradisional dengan LSM-Tree (Log-Structured Merge-tree) yang digunakan pada engine penyimpanan modern untuk write-heavy workloads.",
			Category:    "Database",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusDraft,
		},
		{
			Title:       "Catatan Diskusi Arsitektur Monolith Framework PHP Yang Sudah Tidak Digunakan Lagi",
			Content:     "Dokumen ini merupakan catatan diskusi arsitektur sistem lama berarsitektur monolith berbasis PHP yang kini telah tidak digunakan lagi setelah migration penuh ke microservice Golang dan Next.js App Router.",
			Category:    "Arsip",
			CreatedDate: now,
			UpdatedDate: now,
			Status:      domain.StatusThrash,
		},
		{
			Title:       "Draft Riset Caching Layer Menggunakan Memcached Yang Dipindahkan Ke Trash",
			Content:     "Riset mengenai penggunaan Memcached sebagai caching layer temporary yang kini dipindahkan ke kategori sampah (thrash) setelah tim memilih strategi caching yang terintegrasi langsung di memory application layer.",
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
		// Replace old seeds or populate if empty or count < 5
		if count < 5 {
			log.Println("Seeding rich human-like dummy articles to TiDB Cloud database...")
			dummies := r.getDummyArticles()
			for _, article := range dummies {
				if err := r.db.Create(&article).Error; err != nil {
					log.Printf("Error seeding article: %v", err)
				}
			}
			log.Println("Database TiDB Cloud successfully seeded with realistic dummy articles!")
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
