# 🚀 Sharing Vision - Article Microservice (Backend Golang)

Selamat datang di repository Backend **Article Microservice** untuk Technical Test **Sharing Vision**. Microservice ini dibangun dengan fokus utama pada performa tinggi, struktur kode yang bersih (*Clean Architecture*), skala yang mudah dikembangkan (*scalable*), serta validasi data yang ketat.

---

## 🛠️ Tech Stack & Library

* **Bahasa Pemrograman**: [Golang v1.22+](https://golang.org/)
* **Web Framework**: [Gin Framework](github.com/gin-gonic/gin) - Framework HTTP cepat dan ringan.
* **ORM & Database**: [GORM](gorm.io/gorm) + [MySQL Driver](gorm.io/driver/mysql) - Abstraksi database & auto-migration.
* **Validator**: [Go-Playground Validator v10](github.com/go-playground/validator/v10) - Validasi payload JSON presisi.
* **CORS**: `github.com/gin-contrib/cors` - Penanganan Cross-Origin Resource Sharing.

---

## 🏗️ Arsitektur Proyek (Clean Architecture)

Proyek ini menerapkan **Clean Architecture** untuk memisahkan tanggung jawab (*separation of concerns*) di setiap layer:

```text
Backend/
├── config/             # Konfigurasi database MySQL & environment
├── domain/             # Entity model utama (Tabel `posts`)
├── dto/                # Data Transfer Object & aturan validasi struct
├── repository/         # Query database (GORM) + In-Memory Fallback
├── service/            # Business Logic & transformasi data DTO
├── handler/            # Gin Controller & Formatter error validasi HTTP
├── router/             # Definisi route API & Middleware CORS
├── migrations/         # SQL Script Migration Database (Up/Down)
├── go.mod              # Dependency Go Modules
└── main.go             # Entrypoint aplikasi
```

---

## 🗄️ Database Schema & Migration

Tabel utama bernama **`posts`** dengan struktur kolom sebagai berikut:

| Nama Kolom | Tipe Data | Keterangan |
| :--- | :--- | :--- |
| `id` | `INT` | Auto Increment, Primary Key |
| `title` | `VARCHAR(200)` | Judul Artikel |
| `content` | `TEXT` | Isi Lengkap Artikel |
| `category` | `VARCHAR(100)` | Kategori Artikel |
| `created_date` | `TIMESTAMP` | Waktu Pembuatan Artikel |
| `updated_date` | `TIMESTAMP` | Waktu Terakhir Diperbarui |
| `status` | `VARCHAR(100)` | Enum / Status (`publish` \| `draft` \| `thrash`) |

Script migration manual dapat ditemukan di folder `migrations/000001_create_posts_table.up.sql`.

---

## 📡 API Endpoints & Validasi

### Aturan Validasi Input JSON:
- `title`: Wajib diisi (**Required**), minimal **20 karakter**.
- `content`: Wajib diisi (**Required**), minimal **200 karakter**.
- `category`: Wajib diisi (**Required**), minimal **3 karakter**.
- `status`: Wajib diisi (**Required**), pilihan: `"publish"`, `"draft"`, atau `"thrash"`.

### Daftar Endpoints:

| Method | Endpoint | Deskripsi | Request Body Contoh |
| :--- | :--- | :--- | :--- |
| **POST** | `/article/` | Membuat artikel baru | `{"title": "...", "content": "...", "category": "...", "status": "publish"}` |
| **GET** | `/article/:limit/:offset` | Mengambil list artikel (Paging) | - |
| **GET** | `/article/:id` | Mengambil detail artikel berdasarkan ID | - |
| **PATCH / PUT** | `/article/:id` | Mengubah artikel berdasarkan ID | `{"title": "...", "content": "...", "category": "...", "status": "draft"}` |
| **DELETE** | `/article/:id` | Soft delete artikel (ubah status ke `thrash`) | - |

---

## 💻 Cara Menjalankan Aplikasi

### 1. Prasyarat
- Go 1.22 atau yang lebih baru.
- Database MySQL (Opsional, server otomatis menggunakan *in-memory fallback* jika database belum terhubung).

### 2. Konfigurasi Environment (Opsional)
Anda dapat menentukan variabel lingkungan berikut atau menggunakan nilai default:
```bash
export DB_USER=root
export DB_PASSWORD=your_password
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_NAME=article
export PORT=8080
```

### 3. Menjalankan Server
```bash
# Download dependencies
go mod tidy

# Jalankan server
go run main.go
```
Server akan aktif di `http://localhost:8080`.

---

## 📮 Postman Collection

Collection Postman v2.1 siap diimpor. Anda cukup mengambil file JSON Postman Collection yang disertakan pada dokumentasi atau mengimpor file `postman_collection.json` untuk menguji semua endpoint.
