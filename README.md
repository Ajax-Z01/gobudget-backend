# GoBudget - Personal Finance Tracker

GoBudget adalah aplikasi pencatat keuangan pribadi berbasis API yang dibuat dengan Golang dan PostgreSQL. Aplikasi ini memungkinkan pengguna untuk mencatat pemasukan dan pengeluaran serta melihat ringkasan keuangan mereka.

## Fitur
- Menambahkan transaksi (pemasukan dan pengeluaran) dengan kategori
- Mengupdate transaksi
- Menghapus transaksi (soft delete dengan `deleted_at`)
- Melihat daftar transaksi
- Memfilter transaksi berdasarkan kategori, tipe, dan rentang tanggal
- Melihat ringkasan keuangan (total pemasukan, pengeluaran, dan saldo)
- Seeding data untuk testing

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
### 1. Menampilkan Semua Transaksi
```http
GET /transactions
```

### 2. Menampilkan Transaksi dengan Filter
```http
GET /transactions?category_id=1&type=Income&start_date=2025-03-01&end_date=2025-03-03
```

### 3. Menambahkan Transaksi
```http
POST /transactions
```
Body JSON:
```json
{
  "type": "Income",
  "amount": 5000,
  "note": "Gaji bulan ini",
  "category_id": 1
}
```

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

### 4. Mengupdate Transaksi
```http
PUT /transactions/{id}
```
Body JSON:
```json
{
  "amount": 5500,
  "note": "Gaji bulan ini (update)",
  "category_id": 2
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

### 4. Menghapus Transaksi (Soft Delete)
```http
DELETE /transactions/{id}
```

### 5. Ringkasan Keuangan
```http
GET /summary
```
Response:
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
