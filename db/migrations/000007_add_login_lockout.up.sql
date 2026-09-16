-- Migration 000007: kunci akun otomatis setelah 5x gagal login.
-- Kolom ini bisa dipantau/di-reset manual via Beekeeper Studio:
--   failed_login_attempts = jumlah gagal beruntun (reset ke 0 saat sukses)
--   locked_until          = terkunci sampai waktu ini (NULL = tidak terkunci)
-- Buka kunci manual: UPDATE users SET failed_login_attempts = 0, locked_until = NULL WHERE email = '...';

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS failed_login_attempts INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ NULL;
