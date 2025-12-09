
# 📈 FuncPro API - Komoditas dan Analisis Harga

FuncPro adalah layanan *backend* API yang dirancang untuk melacak, mengelola, dan menganalisis data harga berbagai komoditas. Aplikasi ini dibangun menggunakan Go dan GIN *framework* dan siap untuk di-*deploy* menggunakan **Docker**.

-----

## 💻 Teknologi yang Digunakan

  * **Bahasa Pemrograman:** Go (Golang)
  * **Web Framework:** GIN Web Framework
  * **Database:** PostgreSQL
  * **ORM (Object-Relational Mapper):** GORM
  * **Containerization:** **Docker** dan **Docker Compose**

-----

## 📂 Struktur Proyek

Struktur direktori proyek ini telah dioptimalkan untuk *containerization* dan pengembangan terstruktur.

```
.
├── cmd/
│   └── main.go           # Titik masuk utama aplikasi (entry point).
├── internal/
│   └── config/           # Menangani loading konfigurasi aplikasi (environment variables).
├── pkg/                  # Berisi modul-modul bisnis (komoditas, price, dll.).
├── db/                   # Logic koneksi dan inisialisasi database (GORM AutoMigrate).
├── router/               # Definisi semua routes GIN.
├── Dockerfile            # Instruksi untuk membangun image Docker aplikasi Go.
├── docker-compose.yml    # Konfigurasi untuk menjalankan aplikasi dan PostgreSQL.
└── README.md
```

-----

## 🚀 Endpoint API

Semua *endpoint* memiliki *prefix* `/api/v1`.

| Metode | Path | Deskripsi |
| :--- | :--- | :--- |
| **GET** | `/komoditas` | Mengambil daftar semua komoditas. |
| **POST** | `/komoditas` | Membuat komoditas baru. |
| **GET** | `/komoditas/:id` | Mengambil detail komoditas. |
| **PUT** | `/komoditas/:id` | Memperbarui data komoditas. |
| **DELETE** | `/komoditas/:id` | Menghapus komoditas (*soft delete*). |
| **GET** | `/komoditas/:id/stats` | **Analisis:** Mengambil detail komoditas beserta data statistik harga (Avg, Min, Max, Count, Trend). |
| **POST** | `/prices` | Membuat satu data harga baru. |
| **POST** | `/prices/bulk` | Memasukkan banyak data harga sekaligus (*bulk insert*). |
| **GET** | `/prices/komoditas/:komoditas_id` | Mengambil semua data harga untuk ID komoditas tertentu. |
| **GET** | `/prices/komoditas/:komoditas_id/analysis` | **Analisis:** Mengambil data harga mentah untuk analisis historis. |
| **GET** | `/health` | Mengembalikan status OK. |

-----

## 🐳 Deployment (Docker & Docker Compose)

Aplikasi ini dirancang untuk berjalan di dalam *container* Docker, bersamaan dengan *container* PostgreSQL.

### Prasyarat

  * **Docker**
  * **Docker Compose**

### Langkah 1: Konfigurasi Environment

Pastikan Anda memiliki file `.env` di *root* proyek. Konfigurasi ini akan digunakan oleh aplikasi Go **dan** *container* PostgreSQL.

```env
# Konfigurasi Database (Digunakan oleh Go dan Docker Compose)
DB_HOST=postgres_db  # Nama service Docker Compose
DB_USER=postgres
DB_PASSWORD=secretpassword
DB_NAME=FUNCPRO
DB_PORT=5432

# Konfigurasi Aplikasi (Digunakan oleh Go)
SERVER_PORT=8080
```

### Langkah 2: Build dan Run Menggunakan Docker Compose

Jalankan perintah berikut di *root* direktori proyek Anda.

```bash
# Membangun image aplikasi Go dan memulai semua services (app dan db)
docker-compose up --build
```

### Langkah 3: Verifikasi

Setelah *container* berjalan, aplikasi API Anda akan tersedia di port `8080` pada host lokal Anda.

  * **Akses API:** `http://localhost:8080/api/v1/health`

### Detail Docker Compose (`docker-compose.yml`)

File `docker-compose.yml` mendefinisikan dua *service*:

1.  **`app` (Go Application):**
      * Membangun dari `Dockerfile` di direktori saat ini.
      * Menggunakan *environment variables* dari file `.env`.
      * *Port* **8080** diekspos ke host lokal.
2.  **`postgres_db` (Database):**
      * Menggunakan *image* resmi `postgres:latest`.
      * Menggunakan *environment variables* (`POSTGRES_USER`, `POSTGRES_DB`, dll.) dari file `.env` untuk inisialisasi *database*.
      * Volume digunakan untuk *persistence* data.

-----

## 🛠️ Pengembangan Lokal (Tanpa Docker)

Jika Anda ingin menjalankan aplikasi secara langsung di mesin lokal:

1.  Pastikan PostgreSQL Anda berjalan dan *database* `FUNCPRO` sudah dibuat.
2.  Pastikan file `.env` Anda sudah dikonfigurasi dengan `DB_HOST` yang mengarah ke `localhost`.
3.  Jalankan aplikasi:
    ```bash
    go run cmd/main.go
    ```
