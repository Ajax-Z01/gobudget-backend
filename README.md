# GoBudget - Personal Finance Tracker

GoBudget adalah aplikasi pencatat keuangan pribadi berbasis API yang dibuat dengan Golang dan PostgreSQL. Aplikasi ini memungkinkan pengguna untuk mencatat pemasukan dan pengeluaran serta melihat ringkasan keuangan mereka.

## Fitur Utama
- **Menambah Transaksi**: Pemasukan dan pengeluaran dapat dicatat dengan mudah.
- **Melihat Daftar Transaksi**: Semua transaksi yang belum dihapus dapat dilihat.
- **Soft Delete Transaksi**: Transaksi tidak dihapus dari database, tetapi ditandai dengan `deleted_at`.
- **Restore Transaksi**: Transaksi yang sudah dihapus dapat dipulihkan.
- **Ringkasan Keuangan**: Menampilkan total pemasukan, pengeluaran, dan saldo saat ini.

## Teknologi yang Digunakan
- **Golang** sebagai backend utama
- **Gin Gonic** sebagai framework HTTP
- **GORM** sebagai ORM untuk PostgreSQL
- **PostgreSQL** sebagai database utama

## Instalasi dan Menjalankan Aplikasi

### 1. Clone Repository
```sh
git clone https://github.com/Ajax-Z01/gobudget.git
cd gobudget
```

### 2. Konfigurasi Database
Pastikan PostgreSQL sudah berjalan dan buat database dengan nama `gobudget`. Sesuaikan konfigurasi di file `config.go`:
```go
DB_USER = "postgres"
DB_PASSWORD = "yourpassword"
DB_NAME = "gobudget"
DB_HOST = "localhost"
DB_PORT = "5432"
```

### 3. Jalankan Aplikasi
```sh
go run main.go
```

## Endpoint API

### 1. Menambah Transaksi
- **Endpoint:** `POST /transactions`
- **Body JSON:**
  ```json
  {
    "type": "Income",
    "amount": 5000,
    "note": "Gaji bulan ini"
  }
  ```

### 2. Melihat Semua Transaksi
- **Endpoint:** `GET /transactions`
- **Response:**
  ```json
  [
    {
      "id": 1,
      "type": "Income",
      "amount": 5000,
      "note": "Gaji bulan ini",
      "deleted_at": null
    }
  ]
  ```

### 3. Soft Delete Transaksi
- **Endpoint:** `PUT /transactions/delete/:id`
- **Deskripsi:** Menandai transaksi sebagai dihapus dengan mengisi `deleted_at`.

### 4. Restore Transaksi
- **Endpoint:** `PUT /transactions/restore/:id`
- **Deskripsi:** Mengembalikan transaksi yang sudah dihapus dengan mengosongkan `deleted_at`.

### 5. Melihat Ringkasan Keuangan
- **Endpoint:** `GET /summary`
- **Response:**
  ```json
  {
    "balance": 2000,
    "total_expense": 3000,
    "total_income": 5000
  }
  ```

## Kontribusi
Silakan buat pull request atau issue jika ingin berkontribusi pada proyek ini.

## Lisensi
Proyek ini menggunakan lisensi MIT.
