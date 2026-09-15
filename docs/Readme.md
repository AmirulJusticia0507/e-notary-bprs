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

## Quick Start (Local Setup)

### Backend (Go)
```bash
cd backend
go mod download
go run main.go
```

### Frontend (Vue 3 + Tailwind)
```bash
cd frontend
npm install
npm run dev
```

### Database (PostgreSQL Migration)
```bash
golang-migrate -path db/migrations -database "postgres://user:pass@localhost:5432/bprs_enotary?sslmode=disable" up
```