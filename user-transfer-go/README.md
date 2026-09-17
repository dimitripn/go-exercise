# User Transfer API

REST API sederhana untuk manajemen user dan transfer saldo antar user,
dibangun dengan Go + framework [Gin](https://gin-gonic.com/), dengan
penyimpanan data di **PostgreSQL**.

Struktur project dipisah per layer (repository, service, handler, routes)
mengikuti pola clean/layered architecture yang umum dipakai di Go, supaya
business logic tidak tercampur dengan kode HTTP maupun kode akses database.

## Cara saya membangun aplikasi ini (alur berpikir / starting point)

Saya mulai dari spesifikasi endpoint dan skema data yang diberikan, lalu
menyusunnya dengan pendekatan "dari data ke luar" (inside-out), seperti
biasanya alur development REST API di Go:

1. **Model / Entity** (`internal/models`)
   Definisikan dulu bentuk data inti: `User` (id, name, age, balance) dan
   `Transfer` (id, user_id, target_user_id, nominal, created_date). Ini jadi
   kontrak yang dipakai semua layer di atasnya. Request body (DTO) dipisah
   dari entity (`CreateUserRequest`, `CreateTransferRequest`) supaya validasi
   input tidak bocor ke struktur data internal.

2. **Migration SQL** (`migrations/001_init.sql`)
   Skema tabel `users` dan `transfers` didefinisikan lebih dulu sebelum kode
   Go ditulis, memakai `SERIAL` bawaan Postgres untuk primary key (auto
   increment, tidak perlu generate id manual), plus foreign key
   `transfers.user_id`/`target_user_id -> users.id` dan index untuk query
   riwayat transfer.

3. **Config & Database connection** (`internal/config`, `internal/database`)
   - `internal/config` membaca kredensial koneksi dari environment variable
     (opsional lewat file `.env` saat development, via `godotenv`).
   - `internal/database` membuka connection pool ke Postgres memakai driver
     `pgx` (lewat `database/sql` standar) dan melakukan `Ping` supaya error
     konfigurasi/koneksi ketahuan langsung saat start, bukan saat request
     pertama masuk.

4. **Repository** (`internal/repository`)
   Layer akses data paling bawah, isinya query SQL murni (`database/sql`)
   ke tabel `users` dan `transfers`. Layer ini hanya tahu cara CRUD data,
   tidak tahu apa itu "aturan bisnis". Semua query pakai parameterized
   query (`$1`, `$2`, ...) untuk mencegah SQL injection.

5. **Service** (`internal/service`)
   Di sinilah business rule diletakkan, terpisah dari HTTP dan storage:
   - user/target harus ada sebelum transfer diproses
   - tidak boleh transfer ke diri sendiri
   - nominal harus > 0
   - saldo pengirim harus mencukupi
   - proses baca saldo -> validasi -> update saldo -> catat transfer
     dibungkus dalam satu **database transaction** (`sql.Tx`) dengan
     `SELECT ... FOR UPDATE` (row locking) supaya atomic dan aman dari race
     condition saat ada beberapa transfer berjalan bersamaan.
   Error dideklarasikan sebagai sentinel error (`ErrUserNotFound`, dll) agar
   handler bisa memetakannya ke HTTP status code yang tepat memakai
   `errors.Is`.

6. **Handler** (`internal/handlers`)
   Adapter HTTP <-> service. Tugasnya: parsing/validasi request (path
   param, JSON body via `ShouldBindJSON` + struct tag `binding`), memanggil
   service, lalu menerjemahkan hasil/error jadi response JSON + status code
   yang sesuai (200/201/204/400/404/500).

7. **Routes** (`internal/routes`)
   Mendaftarkan seluruh path Gin ke handler terkait, dikelompokkan dengan
   `router.Group("/users")` supaya nested route (`/users/:id/transfers`)
   rapi.

8. **main.go**
   Bertugas: load config -> connect database -> dependency wiring
   (repository -> service -> handler -> routes) -> jalankan server Gin.
   Tidak ada business logic di sini.

Pendekatan berlapis ini (repository -> service -> handler) dipilih supaya:
- business logic (aturan transfer) tidak tercampur dengan kode HTTP/SQL,
  sehingga bisa di-unit-test tanpa perlu HTTP server maupun database asli.
- setiap layer punya tanggung jawab tunggal (single responsibility),
  memudahkan orang lain membaca kode dan menambah fitur.
- gampang ganti driver/database lain di masa depan: cukup ganti isi
  `internal/database` dan `internal/repository`, layer service & handler
  tidak perlu berubah.

## Struktur folder

```
.
├── main.go                          # entrypoint: load config, connect DB, wiring, start server
├── migrations/
│   └── 001_init.sql                 # skema tabel users & transfers
├── .env.example                     # contoh konfigurasi koneksi DB
├── internal/
│   ├── config/
│   │   └── config.go                # baca konfigurasi dari env var / .env
│   ├── database/
│   │   └── database.go              # koneksi ke Postgres (database/sql + pgx)
│   ├── models/
│   │   ├── user.go                  # entity User
│   │   ├── transfer.go              # entity Transfer
│   │   └── dto.go                   # request/response DTO
│   ├── repository/
│   │   ├── user_repository.go       # query SQL untuk tabel users
│   │   └── transfer_repository.go   # query SQL untuk tabel transfers
│   ├── service/
│   │   ├── errors.go                # sentinel error business logic
│   │   ├── user_service.go          # business logic user
│   │   └── transfer_service.go      # business logic transfer (validasi saldo, transaction, dll)
│   ├── handlers/
│   │   ├── user_handler.go          # HTTP handler /users
│   │   └── transfer_handler.go      # HTTP handler /users/:id/transfers
│   └── routes/
│       └── routes.go                # pendaftaran semua route
├── go.mod
└── go.sum
```

## Persiapan database (Postgres)

1. Pastikan Postgres sudah jalan (lokal atau container), lalu buat database:
   ```bash
   createdb user_transfer
   ```
2. Jalankan migration untuk membuat tabel:
   ```bash
   psql -d user_transfer -f migrations/001_init.sql
   ```
3. Salin `.env.example` menjadi `.env` dan sesuaikan kredensialnya:
   ```bash
   cp .env.example .env
   ```
   Isi `.env`:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=user_transfer
   DB_SSLMODE=disable
   APP_PORT=8080
   ```
   File `.env` sudah di-ignore lewat `.gitignore` supaya kredensial tidak
   ikut ter-commit. `internal/config` akan otomatis membaca file ini saat
   development; di production cukup set environment variable yang sama
   langsung di OS/container, tidak perlu file `.env`.

## Menjalankan aplikasi

```bash
go run main.go
```

Aplikasi akan gagal start (fail fast) dengan pesan error yang jelas kalau
koneksi ke database gagal — jadi masalah konfigurasi ketahuan sejak awal,
bukan saat ada request masuk.

Server berjalan di `http://localhost:8080` (atau sesuai `APP_PORT`).

## Endpoint API

### 1. Get semua user
```
GET /users
```
Response `200`:
```json
[
  { "id": 1, "name": "Budi", "age": 25, "balance": 100000 }
]
```

### 2. Get user berdasarkan ID
```
GET /users/{user_id}
```
- `200` — data user
- `404` — `{"error": "user tidak ditemukan"}`

### 3. Buat user baru
```
POST /users
Content-Type: application/json

{
  "name": "Budi",
  "age": 25,
  "balance": 100000
}
```
- `201` — user berhasil dibuat, `id` di-generate otomatis oleh Postgres (SERIAL)
- `400` — validasi gagal (name kosong, age <= 0, balance negatif)

### 4. Hapus user
```
DELETE /users/{user_id}
```
- `204` — berhasil dihapus (tanpa body). Seluruh transfer terkait user ini
  ikut terhapus otomatis (foreign key `ON DELETE CASCADE`).
- `404` — user tidak ditemukan

### 5. Get riwayat transfer milik user
```
GET /users/{user_id}/transfers
```
Mengembalikan seluruh transfer yang melibatkan user tersebut, baik sebagai
pengirim maupun penerima, diurutkan dari yang terbaru. Field `direction`
menandai `"OUT"` (mengirim) atau `"IN"` (menerima).
- `200`:
```json
[
  {
    "id": 1,
    "user_id": 1,
    "target_user_id": 2,
    "nominal": 20000,
    "created_date": "2026-09-17T13:42:48+07:00",
    "direction": "OUT"
  }
]
```
- `404` — user tidak ditemukan

### 6. Buat transfer baru
```
POST /users/{user_id}/transfers
Content-Type: application/json

{
  "target_user_id": 2,
  "nominal": 20000
}
```
`user_id` di path adalah pengirim (sumber saldo dikurangi),
`target_user_id` di body adalah penerima (saldo ditambah).

- `201` — transfer berhasil, saldo kedua user langsung ter-update secara
  atomic dalam satu database transaction
- `400` — nominal <= 0 / transfer ke diri sendiri / saldo tidak mencukupi
- `404` — user pengirim atau target tidak ditemukan

## Contoh alur pengujian manual (curl)

```bash
# buat dua user
curl -X POST localhost:8080/users -H 'Content-Type: application/json' \
  -d '{"name":"Budi","age":25,"balance":100000}'
curl -X POST localhost:8080/users -H 'Content-Type: application/json' \
  -d '{"name":"Siti","age":30,"balance":5000}'

# transfer dari user 1 ke user 2
curl -X POST localhost:8080/users/1/transfers -H 'Content-Type: application/json' \
  -d '{"target_user_id":2,"nominal":20000}'

# lihat riwayat transfer user 1
curl localhost:8080/users/1/transfers

# lihat saldo terbaru semua user
curl localhost:8080/users
```

Seluruh skenario di atas (termasuk error case: user/target tidak ditemukan,
saldo kurang, transfer ke diri sendiri, delete cascade) sudah diuji manual
langsung terhadap Postgres asli selama development dan berjalan sesuai
status code yang didokumentasikan.

## Kenapa pakai `SERIAL` bukan UUID

Sesuai request: id dibuat sesederhana mungkin memakai `SERIAL` bawaan
Postgres (auto-increment integer, mulai dari 1). Tidak perlu generate ID di
sisi aplikasi — cukup `INSERT ... RETURNING id` dan Postgres yang urus.

## Catatan konkurensi / integritas data

Proses transfer (`POST /users/{id}/transfers`) dijalankan dalam satu
database transaction dengan `SELECT ... FOR UPDATE` untuk mengunci baris
user pengirim & penerima sebelum membaca saldo. Ini mencegah dua transfer
yang berjalan bersamaan saling menimpa saldo (race condition) atau
menyebabkan saldo minus. Baris dikunci berurutan berdasarkan id (id lebih
kecil dulu) untuk menghindari deadlock saat ada transfer berlawanan arah
(A->B dan B->A) yang berjalan di saat bersamaan.

## Catatan pengembangan lanjutan

- Tinggal tambah file `migrations/002_xxx.sql` untuk perubahan skema
  berikutnya (bisa dirapikan pakai tool migration seperti `golang-migrate`
  kalau project makin besar).
- Kolom `age` dan `balance` sudah punya `CHECK` constraint di level
  database (`age > 0`, `balance >= 0`, `nominal > 0`) sebagai lapisan
  pertahanan tambahan selain validasi di service/handler.
