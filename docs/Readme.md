# e-Notary BPRS Documentation

Sistem Manajemen Legalitas Akad Pembiayaan Syariah & Order Notaris Rekanan.

## Documentation Index

- [Product Goals](goals.md) — Core objectives & tech stack
- [Dev Guidelines](guidelines.md) — Frontend & backend engineering rules
- [UI Styles](styles.md) — Tailwind CSS palette & component styles
- [Design System](DESIGN.md) — Colors, typography, layout & component specs
- [Table Specs](tables.md) — High-density table patterns
- [Form Specs](forms.md) — Form architecture & validation rules
- [System Structures](structures.md) — Folder layout, architecture & DB schema
- [Setup Guide](SETUP.md) — Migrasi DB, konfigurasi env, verifikasi backend

## Quick Start (Local Setup)

### 1. Database (PostgreSQL Migration)

```bash
createdb -U postgres bprs_enotary
migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/bprs_enotary?sslmode=disable" up
```

Migrasi: `db/migrations/000001_create_users_notaries` → `000002_create_financing_collateral` → `000003_create_legal_orders_documents` → `000004_seed_dev` → `000005_add_financing_nik_unique` (sinkronisasi CBS memerlukan NIK unik).
Detail & troubleshooting: [SETUP.md](SETUP.md).

### 2. Backend (Go)

```bash
cd backend
cp .env.example .env
go mod download
go run ./cmd/api
```

Verifikasi: `gofmt -l .`, `go vet ./...`, `go test ./...`, `go build ./...` (dari `backend/`).

### 3. Frontend (Vue 3 + Tailwind)

```bash
cd frontend
npm install
npm run dev
```
