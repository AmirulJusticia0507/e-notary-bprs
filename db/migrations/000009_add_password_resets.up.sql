-- Migration 000009: token reset password.
-- Alur MVP (belum ada SMTP): nasabah minta reset via /auth/forgot-password,
-- admin melihat token via GET /users/password-resets lalu relay ke nasabah
-- via WA/telepon, nasabah set password baru via /auth/reset-password.
-- token_hash = SHA256 hex dari token; hanya hash yang disimpan.

CREATE TABLE IF NOT EXISTS password_resets (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash CHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_resets_user ON password_resets (user_id);
CREATE INDEX IF NOT EXISTS idx_resets_expiry ON password_resets (expires_at);
