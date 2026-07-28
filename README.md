# 🚀 Sharing Vision - Article Microservice (Backend Golang)

Selamat datang di repository Backend **Article Microservice** untuk Technical Test **Sharing Vision**. Microservice ini dibangun dengan fokus utama pada performa tinggi, struktur kode yang bersih (*Clean Architecture*), skala yang mudah dikembangkan (*scalable*), serta terhubung langsung ke **TiDB Cloud Database**.

---

## 🛠️ Tech Stack & Library

* **Bahasa Pemrograman**: [Golang v1.22+](https://golang.org/)
* **Web Framework**: [Gin Framework](github.com/gin-gonic/gin) - Framework HTTP cepat dan ringan.
* **ORM & Database**: [GORM](gorm.io/gorm) + [MySQL Driver](gorm.io/driver/mysql) - Terhubung ke **TiDB Cloud**.
* **Validator**: [Go-Playground Validator v10](github.com/go-playground/validator/v10) - Validasi payload JSON presisi.
* **CORS**: `github.com/gin-contrib/cors` - Penanganan Cross-Origin Resource Sharing.
* **Deployment & Container**: [Docker](https://www.docker.com/) & [Render.com](https://render.com/)

---

## 🏗️ Arsitektur Proyek (Clean Architecture)

Proyek ini menerapkan **Clean Architecture** untuk memisahkan tanggung jawab (*separation of concerns*) di setiap layer:

```text
Backend/
├── config/             # Konfigurasi database TiDB Cloud & environment
├── domain/             # Entity model utama (Tabel `posts`)
├── dto/                # Data Transfer Object & aturan validasi struct
├── repository/         # Query database (GORM) + TiDB Cloud Seeder
├── service/            # Business Logic & transformasi data DTO
├── handler/            # Gin Controller & Formatter error validasi HTTP
├── router/             # Definisi route API & Middleware CORS
├── migrations/         # SQL Script Migration Database (Up/Down)
├── Dockerfile          # Multi-stage Docker build untuk Render
├── render.yaml         # Blueprint Konfigurasi Deployment Render
├── go.mod              # Dependency Go Modules
└── main.go             # Entrypoint aplikasi
```

---

## ☁️ Tutorial Deployment ke Render (Render.com)

Ada 2 cara mudah untuk me-deploy backend Golang ini ke Render:

### Opsi A: Deploy Menggunakan Dockerfile (Sangat Direkomendasikan)
1. Login ke akun [Render.com](https://dashboard.render.com/).
2. Klik tombol **New +** → Pilih **Web Service**.
3. Hubungkan ke akun GitHub Anda dan pilih repository **`Backend_Go_SVI`** (`https://github.com/rusdifn08/Backend_Go_SVI.git`).
4. Konfigurasikan Service:
   - **Name**: `sharing-vision-backend` (atau nama pilihan Anda)
   - **Region**: `Singapore` (Paling dekat dengan TiDB Cloud AWS ap-southeast-1)
   - **Language / Runtime**: `Docker`
   - **Dockerfile Path**: `Dockerfile`
   - **Instance Type**: `Free`
5. Tambahkan **Environment Variables** di bagian Advanced Settings:
   - `DB_USER` = `4TYLPkigyGPqtu5.root`
   - `DB_PASSWORD` = `eU91Td4D7tdl9YEA`
   - `DB_HOST` = `gateway01.ap-southeast-1.prod.aws.tidbcloud.com`
   - `DB_PORT` = `4000`
   - `DB_NAME` = `test`
6. Klik **Create Web Service**. Render akan otomatis mem-build Docker container dan mendeploy aplikasi Anda!

### Opsi B: Deploy Tanpa Docker (Native Go Web Service)
1. Pilih **Web Service** baru dari GitHub repo `Backend_Go_SVI`.
2. Pilih **Runtime**: `Go`.
3. Set **Build Command**: `go build -o main .`
4. Set **Start Command**: `./main`
5. Masukkan Environment Variables TiDB Cloud yang sama seperti di atas.
6. Klik **Create Web Service**.

---

## 🗄️ Database Schema & Migration

Tabel utama bernama **`posts`** pada TiDB Cloud:

| Nama Kolom | Tipe Data | Keterangan |
| :--- | :--- | :--- |
| `id` | `INT` | Auto Increment, Primary Key |
| `title` | `VARCHAR(200)` | Judul Artikel |
| `content` | `TEXT` | Isi Lengkap Artikel |
| `category` | `VARCHAR(100)` | Kategori Artikel |
| `created_date` | `TIMESTAMP` | Waktu Pembuatan Artikel |
| `updated_date` | `TIMESTAMP` | Waktu Terakhir Diperbarui |
| `status` | `VARCHAR(100)` | Status artikel (`publish` \| `draft` \| `thrash`) |

---

## 📡 API Endpoints & Validasi

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| **POST** | `/article/` | Membuat artikel baru (Validasi: title>=20, content>=200, category>=3, status in publish/draft/thrash) |
| **GET** | `/article/:limit/:offset` | Mengambil list artikel (Paging) |
| **GET** | `/article/:id` | Mengambil detail artikel berdasarkan ID |
| **PATCH / PUT** | `/article/:id` | Mengubah artikel berdasarkan ID |
| **DELETE** | `/article/:id` | Soft delete artikel (ubah status ke `thrash`) |
