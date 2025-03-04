# GoBudget

GoBudget adalah aplikasi pencatatan keuangan pribadi berbasis REST API yang dibuat menggunakan Golang dan PostgreSQL.

## Fitur
- Registrasi dan login pengguna
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
JWT_SECRET=your_jwt_secret
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
### 1. Registrasi Pengguna
```http
POST /register
```
Body JSON:
```json
{
  "name": "Test User",
  "email": "test@example.com",
  "password": "password123"
}
```

### 2. Login Pengguna
```http
POST /login
```
Body JSON:
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```
Response:
```json
{
  "token": "your_jwt_token"
}
```

### 3. Menampilkan Semua Transaksi
```http
GET /transactions
Authorization: Bearer your_jwt_token
```

### 4. Menampilkan Transaksi dengan Filter
```http
GET /transactions?category_id=1&type=Income&start_date=2025-03-01&end_date=2025-03-03
Authorization: Bearer your_jwt_token
```

### 5. Menambahkan Transaksi
```http
POST /transactions
Authorization: Bearer your_jwt_token
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

### 6. Mengupdate Transaksi
```http
PUT /transactions/{id}
Authorization: Bearer your_jwt_token
```
Body JSON:
```json
{
  "amount": 5500,
  "note": "Gaji bulan ini (update)",
  "category_id": 2
}
```

### 7. Menghapus Transaksi (Soft Delete)
```http
DELETE /transactions/{id}
Authorization: Bearer your_jwt_token
```

### 8. Ringkasan Keuangan
```http
GET /summary
Authorization: Bearer your_jwt_token
```
Response:
```json
{
  "balance": 2000,
  "total_expense": 3000,
  "total_income": 5000
}
```

### 9. Menambahkan Kategori
```http
POST /categories
Authorization: Bearer your_jwt_token
```
Body JSON:
```json
{
  "name": "Makanan"
}
```

### 10. Menampilkan Semua Kategori
```http
GET /categories
Authorization: Bearer your_jwt_token
```

## Seeding Data
Untuk mengisi database dengan data awal untuk testing, jalankan:
```sh
go run seeder.go
```

## Lisensi
Proyek ini menggunakan lisensi MIT.

