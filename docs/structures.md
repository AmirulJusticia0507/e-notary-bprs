# System & Project Architecture Structures

Dokumen ini mendefinisikan struktur folder, pola arsitektur software, serta hirarki direktori untuk sistem e-Notaris BPRS.

---

## Frontend Structure (Vue 3 + Tailwind CSS)

Menggunakan **Composition API** (`<script setup>`) dengan pemisahan tegas antara komponen atomic/UI, layout, state (Pinia), dan API services.

```text
frontend/
├── public/
│   └── favicon.ico
├── src/
│   ├── assets/              # Static assets (images, icons, global tailwind.css)
│   ├── components/          # Reusable UI Components
│   │   ├── common/          # Button, Modal, Badge, LoadingSpinner
│   │   ├── forms/           # FormInputs, FileUploader, CurrencyInput
│   │   └── tables/          # DataTable, TablePagination, StatusBadge
│   ├── composables/         # Custom Vue Composables (useAuth, useLegalOrder, useFormValidation)
│   ├── layouts/             # Layout Wrapper (AppLayout, AuthLayout, Sidebar)
│   ├── router/              # Vue Router navigation guards & route definitions
│   ├── services/            # Axios API Client & Endpoint Services
│   │   ├── api.js           # Axios base instance with JWT Interceptors
│   │   ├── authService.js
│   │   └── legalOrderService.js
│   ├── stores/              # Pinia Global State Management (userStore, orderStore)
│   ├── views/               # Page Components (Views/Screens)
│   │   ├── auth/            # LoginView.vue
│   │   ├── dashboard/       # DashboardOverview.vue
│   │   └── orders/          # OrderListView.vue, OrderCreateWizard.vue, OrderDetailView.vue
│   ├── App.vue
│   └── main.js
├── tailwind.config.js
├── vite.config.js
└── package.json
```

---

## Backend Structure (Golang - Clean Architecture)

Menggunakan pola Clean / Hexagonal Architecture untuk memisahkan logika bisnis dari HTTP transport (Gin) dan database driver (pgx).

```text
backend/
├── cmd/
│   └── api/
│       └── main.go          # Entry point aplikasi Go
├── internal/
│   ├── config/              # Environment variables & database config (godotenv)
│   ├── domain/              # Core Domain Entities (Pure Go Structs)
│   │   └── entities.go
│   ├── repository/          # Data Access Layer (PostgreSQL Queries & Transactions)
│   │   ├── interfaces.go
│   │   └── postgres/
│   │       ├── user.go
│   │       ├── notary.go
│   │       ├── financing.go
│   │       ├── collateral.go
│   │       ├── legal_order.go
│   │       ├── legal_order_log.go
│   │       └── legal_document.go
│   ├── usecase/             # Business Logic Layer (SLA Rules, Auto-drafting, Verification)
│   │   └── usecase.go
│   └── delivery/
│       └── http/            # Transport Layer (Gin Handlers & Route Setup)
│           ├── middleware/  # JWT Auth
│           ├── auth_handler.go
│           ├── notary_handler.go
│           ├── financing_handler.go
│           ├── order_handler.go
│           ├── document_handler.go
│           └── utils.go
├── pkg/                     # Utility / Shared Helpers
│   ├── auth/                # JWT & bcrypt
│   ├── hash/                # SHA-256
│   ├── db/                  # Database connection (pgx)
│   ├── pdf/                 # PDF generator & hash calculator (SHA-256)
│   ├── esign/               # Integration Client (Privy / PERURI API)
│   ├── emeterai/            # Integration Client (PERURI e-Meterai)
│   └── response/            # Standard JSON Response Wrappers
├── cmd/api/main.go          # Entry point (duplicate for clarity)
├── go.mod
└── go.sum
```

---

## Database Migrations (golang-migrate)

Migrations live at the repository root:

```
db/
└── migrations/
    ├── 000001_create_users_notaries.up.sql
    ├── 000001_create_users_notaries.down.sql
    ├── 000002_create_financing_collateral.up.sql
    ├── 000002_create_financing_collateral.down.sql
    ├── 000003_create_legal_orders_documents.up.sql
    ├── 000003_create_legal_orders_documents.down.sql
    ├── 000004_seed_dev.up.sql
    ├── 000004_seed_dev.down.sql
    ├── 000005_add_financing_nik_unique.up.sql
    └── 000005_add_financing_nik_unique.down.sql
```

---

## Database Schema Structure (PostgreSQL)

Ringkasan entitas tabel relasional dalam database `bprs_enotary`:

```text
bprs_enotary (DB)
├── users                  (Staff BPRS, Legal Officers, Admins)
├── notaries               (Data Notaris/PPAT Rekanan & Wilayah Kerja)
├── financing_applications (Data Pembiayaan Sync dari CBS/LOS)
├── collaterals            (Objek Jaminan: SHM, SHGB, BPKB, Deposito)
├── legal_orders           (Order Legalitas & Status SLA Workflows)
├── legal_order_logs       (Audit Trail Waktu & Log Perubahan Status)
└── legal_documents        (PDF Hash SHA-256, e-Meterai SN, e-Sign Status)
```

**Key DB Conventions:**

- **Primary Key:** `id BIGINT GENERATED ALWAYS AS IDENTITY` atau `UUIDv4`.
- **Timestamps:** Selalu sertakan `created_at TIMESTAMPTZ` dan `updated_at TIMESTAMPTZ` (WIB UTC+7).
- **Soft Deletes:** `deleted_at TIMESTAMPTZ NULL` pada tabel master data.
- **Foreign Keys:** Dilengkapi klausul `ON DELETE RESTRICT` untuk integritas audit legal.