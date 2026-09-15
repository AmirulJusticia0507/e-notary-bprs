
---
### 📂 `guidelines.md`
```markdown
# 📐 Development & Engineering Guidelines

## Frontend (Vue 3)
* Selalu gunakan `<script setup>` & Composition API.
* Pisahkan *UI Components* (Atomic) dan *View/Page Components*.
* Axios Interceptor untuk injeksi JWT Token dan penanganan error terpusat (`401 Unauthenticated`, `403 Forbidden`).

## Backend (Golang)
* **Architecture:** Clean Architecture / Hexagonal Pattern (Handler/Controller -> Usecase/Service -> Repository -> Postgres).
* **Database Driver:** Postgres Native via `pgx/v5` atau GORM.
* **Concurrency:** Gunakan Go Goroutines untuk background worker (misal: pendaftaran HT-el ke BPN / kirim e-Meterai API).
* **Transactions:** Semua operasi mutasi (Create Order + Audit Log) WAJIB dalam `tx` (PostgreSQL Transaction).

## Database (PostgreSQL)
* Primary Key menggunakan `BIGINT` atau `UUIDv4`.
* Manfaatkan tipe data native Postgres: `JSONB` (untuk metadata fleksibel), `TIMESTAMPTZ` (waktu dengan timezone WIB), dan `ENUM`.
📂 Readme.md
Markdown
# 📜 e-Notary BPRS Core System

Sistem Manajemen Legalitas Akad Pembiayaan Syariah & Order Notaris Rekanan.

## 📁 Documentation Index
- [🎯 Product Goals](goals.md)
- [📐 Dev Guidelines](guidelines.md)
- [🎨 UI Styles](styles.md)
- [📊 Table Specs](tables.md)
- [📝 Form Specs](forms.md)

## 🚀 Quick Start (Local Setup)

### Backend (Go)
```bash
cd backend
go mod download
go run main.go
Frontend (Vue 3 + Tailwind)
Bash
cd frontend
npm install
npm run dev
Database (PostgreSQL Migration)
Bash
golang-migrate -path db/migrations -database "postgres://user:pass@localhost:5432/bprs_enotary?sslmode=disable" up
---
<Elicitations message="Mau lanjut bikin starter code-nya?">    </Elicitations>
