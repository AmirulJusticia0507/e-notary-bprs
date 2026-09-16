# Data & Integrations

Dokumen ini mencatat batas MVP e-Notary BPRS dan sumber data eksternal yang perlu diperlakukan sebagai integrasi resmi, bukan asumsi data lokal.

## AHU Notaris

Halaman `https://ahu.go.id/formasiNotaris` adalah sumber informasi formasi wilayah jabatan notaris yang tersedia, bukan direktori nama notaris lengkap.

Data yang terlihat dari halaman tersebut lebih cocok untuk:

- Nama wilayah
- Formasi jabatan notaris per tahun
- Jumlah permohonan
- Formasi tersedia saat ini

Data tersebut tidak menggantikan master data notaris karena tidak berisi detail operasional seperti:

- Nama notaris rekanan
- Alamat kantor
- Wilayah kerja detail
- Nomor SK atau identitas praktik
- Kontak PIC
- Status kerja sama dengan BPRS

Untuk MVP, master notaris sebaiknya dikelola manual oleh admin atau diimpor dari CSV/Excel internal. Integrasi AHU hanya dilakukan jika tersedia akses resmi/API atau sumber publik yang memang menyediakan direktori notaris.

## Scope Input MVP

Tanpa integrasi BPN, Pegadaian, CBS/LOS BPRS, dan sistem notaris, aplikasi ini berfungsi sebagai workflow internal untuk:

- Input data pembiayaan
- Input data nasabah
- Input data agunan
- Upload dokumen persyaratan
- Membuat order legal ke notaris
- Tracking status order
- Arsip dokumen dan riwayat perubahan status

Konsekuensinya, sebagian besar data masih diinput manual oleh admin BPRS atau petugas legal.

## Integrasi Yang Belum Otomatis

### BPN

Integrasi BPN diperlukan untuk validasi dan proses terkait agunan tanah/bangunan, misalnya:

- Cek sertifikat
- Validasi NIB atau nomor hak
- Status hak tanggungan
- Roya
- Validasi bidang tanah

### Pegadaian

Integrasi Pegadaian hanya relevan jika ada jaminan atau proses yang bersinggungan dengan data Pegadaian, misalnya validasi agunan tertentu atau status gadai.

### Core BPRS / CBS / LOS

Integrasi dengan sistem inti BPRS diperlukan agar data tidak diketik ulang, termasuk:

- Data nasabah
- Nomor CIF
- Data pembiayaan
- Plafon
- Akad
- Tenor
- Cabang
- Status pembiayaan
- Kolektibilitas, jika diperlukan untuk proses legal

### Sistem Notaris

Integrasi dengan sistem notaris diperlukan untuk otomatisasi lanjutan, seperti:

- Jadwal tanda tangan
- Nomor akta
- Status minuta
- Upload dokumen final
- Invoice dan biaya jasa
- Status pengerjaan dari kantor notaris

## Master Data Minimal

Sebelum integrasi eksternal tersedia, sistem perlu memiliki master data manual berikut:

- Master notaris
- Master cabang BPRS
- Master produk pembiayaan
- Master jenis agunan
- Master template dokumen
- Master status workflow legal/notaris

## Rekomendasi Tahapan

Tahap awal cukup gunakan input manual dan import CSV/Excel untuk master notaris. Setelah workflow internal stabil, integrasi resmi dapat ditambahkan bertahap:

1. CBS/LOS BPRS untuk menarik data nasabah dan pembiayaan.
2. Sistem notaris untuk update status order dan dokumen final.
3. BPN untuk validasi agunan tanah/bangunan.
4. Pegadaian jika proses bisnis memang membutuhkan data Pegadaian.
