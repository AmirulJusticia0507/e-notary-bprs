# Development & Engineering Guidelines

## Frontend (Vue 3)

- Selalu gunakan `<script setup>` & Composition API.
- Pisahkan *UI Components* (Atomic) dan *View/Page Components*.
- Axios Interceptor untuk injeksi JWT Token dan penanganan error terpusat (`401 Unauthenticated`, `403 Forbidden`).

## Backend (Golang)

- **Architecture:** Clean Architecture / Hexagonal Pattern (Handler/Controller -> Usecase/Service -> Repository -> Postgres).
- **Database Driver:** Postgres Native via `pgx/v5` atau GORM.
- **Concurrency:** Gunakan Go Goroutines untuk background worker (misal: pendaftaran HT-el ke BPN / kirim e-Meterai API).
- **Transactions:** Semua operasi mutasi (Create Order + Audit Log) WAJIB dalam `tx` (PostgreSQL Transaction).

## Database (PostgreSQL)

- Primary Key menggunakan `BIGINT` atau `UUIDv4`.
- Manfaatkan tipe data native Postgres: `JSONB` (untuk metadata fleksibel), `TIMESTAMPTZ` (waktu dengan timezone WIB), dan `ENUM`.