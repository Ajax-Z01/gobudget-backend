# GoBudget

GoBudget adalah aplikasi pencatatan keuangan pribadi berbasis REST API yang dibuat menggunakan Golang dan PostgreSQL.

## Fitur
- Menambahkan transaksi (pemasukan dan pengeluaran) dengan kategori
- Mengupdate transaksi
- Menghapus transaksi (soft delete dengan `deleted_at`)
- Melihat daftar transaksi
- Memfilter transaksi berdasarkan kategori, tipe, dan rentang tanggal
- Melihat ringkasan keuangan (total pemasukan, pengeluaran, dan saldo)
- Seeding data untuk testing

## Teknologi yang Digunakan
- Golang
- Gin (framework web)
- GORM (ORM untuk database)
- PostgreSQL

## Instalasi & Penggunaan
### 1. Clone Repository
```sh
git clone https://github.com/Ajax-Z01/gobudget.git
cd gobudget
```

### 2. Konfigurasi Database
Pastikan PostgreSQL sudah terinstall dan jalankan servernya. Buat database dengan nama `gobudget`.

### 3. Buat File `.env`
Buat file `.env` di root proyek dan isi dengan:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=gobudget
```

### 4. Install Dependencies
```sh
go mod tidy
```

### 5. Jalankan Aplikasi
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

### 5. Menghapus Transaksi (Soft Delete)
```http
DELETE /transactions/{id}
```

### 6. Ringkasan Keuangan
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

### 7. Menambahkan Kategori
```http
POST /categories
```
Body JSON:
```json
{
  "name": "Makanan"
}
```

### 8. Menampilkan Semua Kategori
```http
GET /categories
```

## Seeding Data
Untuk mengisi database dengan data awal untuk testing, jalankan:
```sh
go run seed.go
```

## Lisensi
Proyek ini menggunakan lisensi MIT.
