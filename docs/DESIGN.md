
---
version: alpha
name: e-Notary BPRS Design System
description: Sistem desain antarmuka e-Notaris BPRS yang profesional, patuh hukum syariah, dan berkerapatan tinggi (high-density).
colors:
  primary: "#0F172A"
  brand: "#2563EB"
  accent: "#1E3A8A"
  neutral-bg: "#F8FAFC"
  surface: "#FFFFFF"
  border: "#E2E8F0"
  text-main: "#0F172A"
  text-muted: "#64748B"
  status-success: "#16A34A"
  status-warning: "#D97706"
  status-danger: "#DC2626"
typography:
  h1:
    fontFamily: Inter, sans-serif
    fontSize: 1.5rem
    fontWeight: 700
    lineHeight: 1.3
  body:
    fontFamily: Inter, sans-serif
    fontSize: 0.875rem
    lineHeight: 1.5
  code:
    fontFamily: JetBrains Mono, monospace
    fontSize: 0.75rem
    lineHeight: 1.4
rounded:
  sm: 4px
  md: 6px
  lg: 8px
spacing:
  xs: 4px
  sm: 8px
  md: 16px
  lg: 24px
components:
  card:
    backgroundColor: "{colors.surface}"
    borderColor: "{colors.border}"
    rounded: "{rounded.md}"
    padding: "{spacing.md}"
  button-primary:
    backgroundColor: "{colors.brand}"
    textColor: "#FFFFFF"
    rounded: "{rounded.md}"
    padding: "8px 16px"
  badge-success:
    backgroundColor: "#DCFCE7"
    textColor: "{colors.status-success}"
    rounded: "9999px"
    padding: "2px 10px"
---
## Overview

Antarmuka e-Notaris BPRS dirancang untuk mendukung operasional perbankan syariah dan notariat. Karakter tampilan harus terasa formal, sangat aman (trustworthy), memiliki kontras yang jelas, serta berkerapatan data tinggi (high-density) agar memudahkan staf legal melihat status berkas dengan cepat.

## Colors

- Primary (#0F172A): Digunakan untuk Top Navbar, Sidebar header, dan elemen navigasi utama.
- Brand (#2563EB): Warna aksen utama Tailwind (`blue-600`) untuk tombol eksekusi (CTA), link, dan state aktif.
- Neutral Background (#F8FAFC): Canvas dasar dashboard (`slate-50`).
- Status Colors: Green/Emerald (#16A34A) untuk Completed/Approved, Amber (#D97706) untuk Pending/Warning SLA, dan Red (#DC2626) untuk Rejected/Overdue SLA.

## Typography

Menggunakan font `Inter` untuk seluruh teks UI. Hierarki visual diatur ketat melalui ukuran font (`text-xs` hingga `text-xl`) dan `font-weight`, bukan dengan mengganti jenis font. Kolom ID/Order/Hash wajib memakai font monospace (`JetBrains Mono`).

## Layout

Dashboard menggunakan layout Sidebar + Header kaku (fixed/sticky). Konten tabel dan form menggunakan grid 12 kolom Tailwind dengan padding halaman konsisten 24px (`p-6`).

## Elevation & Depth

Menggunakan gaya datar (*flat design*) dengan *border* tegas (`border-slate-200`) untuk tabel dan form. Shadow hanya digunakan secara sangat halus (`shadow-sm`) pada Card dan Modal untuk menjaga kerapatan informasi.

## Shapes

Menggunakan sudut tumpul yang konsisten: `rounded-md` (6px) untuk komponen form/tombol dan `rounded-lg` (8px) untuk kontainer/card. Dilarang mencampur sudut serba bulat (*full rounded*) dengan sudut tajam di halaman yang sama.

## Components

- card: Wadah putih (`bg-white`) dengan border `border-slate-200` untuk membungkus statistik dan grup form.
- button-primary: Tombol aksi utama (`bg-blue-600 hover:bg-blue-700 text-white`).
- table-data: Tabel berkerapatan tinggi dengan header gelap (`bg-slate-900 text-slate-200`) dan baris belang/hover (`hover:bg-slate-50`).

## Do's and Don'ts

- Do: Gunakan token warna status yang konsisten untuk menentukan SLA Notaris.
- Do: Pastikan semua elemen kode order atau hash sertifikat memakai format font monospace.
- Don't: Menggunakan warna-warni yang tidak ada di palet (misal: ungu/pink) pada status pembiayaan.
- Don't: Membuat spacing komponen terlalu renggang yang menghabiskan ruang layar monitor staf.
