# Setup & Development Guide — e-Notary BPRS

Panduan menjalankan backend Go, database PostgreSQL (migrasi), dan frontend secara lokal.

## 1. Prasyarat

- Go 1.26+ (`go version`)
- PostgreSQL 15+ (lokal atau Docker)
- Node.js 20+ + npm (untuk frontend)
- Salah satu migration runner:
  - `golang-migrate` (`migrate`), atau
  - `psql` langsung (file SQL bersifat idempotent parsial — disarankan via `migrate`)

## 2. Struktur relevan

```text
.
├── backend/                 # Go module github.com/e-notary-bprs/backend
│   ├── cmd/api/main.go      # entrypoint (delivery.NewRouter + pkg/db.Open)
│   ├── internal/config/     # env loader (godotenv) + DSN
│   ├── internal/...         # domain, repository/postgres, usecase, delivery/http (Gin)
│   ├── pkg/...              # auth (JWT+bcrypt), response (gin), db, hash, pdf, esign, emeterai
│   └── .env.example         # contoh env backend
├── db/migrations/           # 000001..000005 *.up.sql / *.down.sql
├── docs/                    # spesifikasi + panduan ini
└── .env.example             # contoh env root
```

## 3. Setup database

```bash
# 1. Buat database
createdb -U postgres bprs_enotary
# atau via psql:
# CREATE DATABASE bprs_enotary;

# 2. Jalankan migrasi (dari root repo)
migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/bprs_enotary?sslmode=disable" up

# Rollback bila perlu:
# migrate -path db/migrations -database "..." down 1
```

Urutan migrasi:

| File | Isi |
| ---- | --- |
| `000001_create_users_notaries` | `users`, `notaries` + index |
| `000002_create_financing_collateral` | `financing_applications`, `collaterals` + index |
| `000003_create_legal_orders_documents` | `legal_orders`, `legal_order_logs`, `legal_documents` + index |
| `000004_seed_dev` | seed dev idempotent: admin `admin@bprs.local` / `admin123`, 1 notaris contoh |
| `000005_add_financing_nik_unique` | unique index untuk `financing_applications.customer_nik` agar upsert CBS aman |

Konvensi DB: PK `BIGINT GENERATED ALWAYS AS IDENTITY`, `created_at/updated_at TIMESTAMPTZ DEFAULT NOW()`,
soft-delete `deleted_at NULL`, FK `ON DELETE RESTRICT` (kecuali `legal_order_logs.order_id → CASCADE`).

## 4. Setup backend

```bash
cd backend
cp .env.example .env        # sesuaikan DB_*, JWT_SECRET
go mod download
go run ./cmd/api            # http://0.0.0.0:8080
```

Variabel env (lihat `internal/config/config.go`):

| Key | Default | Keterangan |
| --- | ------- | ---------- |
| `DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME/DB_SSLMODE` | `localhost/5432/postgres/postgres/bprs_enotary/disable` | koneksi pgx stdlib (`pkg/db.Open(cfg.Database)`) |
| `DATABASE_URL` / `POSTGRES_URL` | kosong | DSN database cloud (Neon/Vercel Postgres); mengambil prioritas atas `DB_*` |
| `SERVER_HOST/SERVER_PORT` | `0.0.0.0/8080` | listen address lokal |
| `PORT` | `SERVER_PORT`/8080 | port runtime cloud; Vercel mengisi otomatis |
| `JWT_SECRET/JWT_EXPIRY` | `e-notary-bprs-secret-key/24h` | HS256, durasi `time.ParseDuration` |

## 5. Deployment Vercel backend

Proyek Vercel untuk API harus menggunakan **Root Directory** `backend`. Konfigurasi [`backend/vercel.json`](../backend/vercel.json) memaksa framework preset **Go** dan membangun binary dari `cmd/api/main.go`.

Tambahkan environment variable berikut di Vercel:

- `DATABASE_URL` atau `POSTGRES_URL`: DSN PostgreSQL cloud, misalnya dari Neon/Vercel Postgres. Variabel ini mengambil prioritas atas konfigurasi `DB_*`.
- `JWT_SECRET`: secret minimal 32 karakter.
- `PORT` tidak perlu diisi; Vercel menetapkannya otomatis dan server membacanya sebelum `SERVER_PORT`.

Jangan gunakan `localhost` untuk database di Vercel. Jalankan migration terhadap database cloud sebelum deployment, lalu redeploy setelah environment variable tersedia.

## 6. Verifikasi backend

```bash
cd backend
gofmt -l .                  # harus kosong
go vet ./...
go test ./...
go build ./...
# atau: make fmt vet test build
```

Alternatif Docker (dari root repo):

```bash
cp backend/.env.example backend/.env   # sesuaikan bila perlu
docker compose up --build              # postgres + migrate + api
```

Test yang tersedia: `internal/config` (defaults + DSN + env override),
`pkg/hash` (FromString/FromReader/FromBytes SHA-256), `pkg/auth` (GenerateToken/ParseToken + bcrypt roundtrip),
`internal/usecase` (transisi status order + SLA dengan stub repo, tanpa DB).

## 6. API ringkas

Publik: `POST /api/v1/auth/register`, `POST /api/v1/auth/login` (balikan `{token, user}`).
Terproteksi (`Authorization: Bearer <JWT>`): `/notaries`, `/financings` (+`POST /sync`),
`/orders` (+`/notary/:notaryID`, `/assigned/:assignedTo`, `PATCH /:id/status`, `GET /:id/logs`),
`/documents` (+`/order/:orderID`, `PATCH /:id/sign`).

Aturan status order: `pending → in_progress → completed|rejected`.
`UpdateStatus` menolak transisi invalid (`400`), order hilang (`404`), dan SLA overdue (`403`,
kecuali menuju `completed`).

## 7. Frontend (ringkas)

```bash
cd frontend
npm install
npm run dev
```

Pastikan `VITE_API_BASE_URL=http://localhost:8080/api/v1` dan JWT dikirim via Axios interceptor.

## 8. Troubleshooting

- `missing go.sum entry`: jalankan `go mod tidy` di `backend/`.
- `stdlib.Open undefined`: gunakan `sql.Open("pgx", dsn)` (sudah diperbaiki di `pkg/db/conn.go`).
- `401 missing/invalid token`: cek `JWT_SECRET` sama antara penerbit token dan validator, format header `Bearer <token>`.
- `relation does not exist`: migrasi belum jalan — ulangi langkah 3.
