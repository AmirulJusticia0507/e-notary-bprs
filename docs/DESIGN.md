# e-Notary BPRS Design System

Sistem desain antarmuka e-Notaris BPRS yang profesional, patuh hukum syariah, dan berkerapatan tinggi (high-density).

## Colors

| Token | Hex | Usage |
| --- | --- | --- |
| Primary | `#0F172A` | Top Navbar, Sidebar header, elemen navigasi utama |
| Brand | `#2563EB` | Tombol eksekusi (CTA), link, state aktif |
| Accent | `#1E3A8A` | Aksen sekunder |
| Neutral Background | `#F8FAFC` | Canvas dasar dashboard |
| Surface | `#FFFFFF` | Permukaan kartu & tabel |
| Border | `#E2E8F0` | Border halus |
| Text Main | `#0F172A` | Teks utama |
| Text Muted | `#64748B` | Teks pendukung |
| Success | `#16A34A` | Completed / Approved |
| Warning | `#D97706` | Pending / Warning SLA |
| Danger | `#DC2626` | Rejected / Overdue SLA |

## Typography

Menggunakan font `Inter` untuk seluruh teks UI. Hierarki visual diatur ketat melalui ukuran font (`text-xs` hingga `text-xl`) dan `font-weight`, bukan dengan mengganti jenis font. Kolom ID/Order/Hash wajib memakai font monospace (`JetBrains Mono`).

## Layout

Dashboard menggunakan layout Sidebar + Header kaku (fixed/sticky). Konten tabel dan form menggunakan grid 12 kolom Tailwind dengan padding halaman konsisten 24px (`p-6`).

## Elevation & Depth

Menggunakan gaya datar (*flat design*) dengan *border* tegas (`border-slate-200`) untuk tabel dan form. Shadow hanya digunakan secara sangat halus (`shadow-sm`) pada Card dan Modal untuk menjaga kerapatan informasi.

## Shapes

Menggunakan sudut tumpul yang konsisten: `rounded-md` (6px) untuk komponen form/tombol dan `rounded-lg` (8px) untuk kontainer/card. Dilarang mencampur sudut serba bulat (*full rounded*) dengan sudut tajam di halaman yang sama.

## Components

- **card:** Wadah putih (`bg-white`) dengan border `border-slate-200` untuk membungkus statistik dan grup form.
- **button-primary:** Tombol aksi utama (`bg-blue-600 hover:bg-blue-700 text-white`).
- **table-data:** Tabel berkerapatan tinggi dengan header gelap (`bg-slate-900 text-slate-200`) dan baris hover (`hover:bg-slate-50`).
- **badge-success:** Badge keberhasilan dengan background `#DCFCE7` dan teks `#16A34A`.

## Do's and Don'ts

- **Do:** Gunakan token warna status yang konsisten untuk menentukan SLA Notaris.
- **Do:** Pastikan semua elemen kode order atau hash sertifikat memakai format font monospace.
- **Don't:** Menggunakan warna-warni yang tidak ada di palet (misal: ungu/pink) pada status pembiayaan.
- **Don't:** Membuat spacing komponen terlalu renggang yang menghabiskan ruang layar monitor staf.