
# 🏗️ System & Project Architecture Structures

Dokumen ini mendefinisikan struktur folder, pola arsitektur software, serta hirarki direktori untuk sistem e-Notaris BPRS.

---

## 🎨 1. Frontend Structure (Vue 3 + Tailwind CSS)

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



⚙️ 2. Backend Structure (Golang - Clean Architecture)
Menggunakan pola Clean / Hexagonal Architecture untuk memisahkan logika bisnis dari HTTP transport (Gin/Fiber) dan database driver (pgx/GORM).

Plaintext
backend/
├── cmd/
│   └── api/
│       └── main.go          # Entry point aplikasi Go
├── config/                  # Environment variables & database config (Viper/godotenv)
├── internal/
│   ├── delivery/
│   │   └── http/            # Transport Layer (Gin Handlers & Route Setup)
│   │       ├── middleware/  # JWT Auth, CORS, Request Logger
│   │       ├── auth_handler.go
│   │       ├── notary_handler.go
│   │       └── order_handler.go
│   ├── domain/              # Core Domain Entities & Interfaces (Pure Go Structs)
│   │   ├── notary.go
│   │   ├── order.go
│   │   └── document.go
│   ├── repository/          # Data Access Layer (PostgreSQL Queries & Transactions)
│   │   ├── postgres/
│   │   │   ├── notary_repository.go
│   │   │   └── order_repository.go
│   │   └── interfaces.go
│   └── usecase/             # Business Logic Layer (SLA Rules, Auto-drafting, Verification)
│       ├── notary_usecase.go
│       └── order_usecase.go
├── pkg/                     # Utility / Shared Helpers
│   ├── pdf/                 # PDF generator & hash calculator (SHA-256)
│   ├── esign/               # Integration Client (Privy / PERURI API)
│   └── response/            # Standard JSON Response Wrappers
├── db/
│   └── migrations/          # SQL Migration Files (golang-migrate)
│       ├── 000001_create_users_table.up.sql
│       └── 000001_create_users_table.down.sql
├── go.mod
└── go.sum
🗄️ 3. Database Schema Structure (PostgreSQL)
Ringkasan entitas tabel relasional dalam database bprs_enotary:

Plaintext
bprs_enotary (DB)
 ├── users                  (Staff BPRS, Legal Officers, Admins)
 ├── notaries               (Data Notaris/PPAT Rekanan & Wilayah Kerja)
 ├── financing_applications (Data Pembiayaan Sync dari CBS/LOS)
 ├── collaterals            (Objek Jaminan: SHM, SHGB, BPKB, Deposito)
 ├── legal_orders           (Order Legalitas & Status SLA Workflows)
 ├── legal_order_logs       (Audit Trail Waktu & Log Perubahan Status)
 └── legal_documents        (PDF Hash SHA-256, e-Meterai SN, e-Sign Status)
Key DB Conventions:
Primary Key: id BIGINT GENERATED ALWAYS AS IDENTITY atau UUIDv4.

Timestamps: Selalu sertakan created_at TIMESTAMPTZ dan updated_at TIMESTAMPTZ (WIB UTC+7).

Soft Deletes: deleted_at TIMESTAMPTZ NULL pada tabel master data.

Foreign Keys: Dilengkapi klausul ON DELETE RESTRICT untuk integritas audit legal.

---

<Elicitations message="Semua 7 file markdown spesifikasi (.md) sudah lengkap! Mau langsung buatkan kode starter Go/Vue-nya?">    </Elicitations>
