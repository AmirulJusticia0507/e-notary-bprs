-- Migration 000008: role nasabah + relasi pengajuan ke akun nasabah.
-- Nasabah daftar sendiri (signup), mengajukan pembiayaan, dan memantau
-- statusnya di dashboard nasabah. user_id NULL = data lama / dari CBS.

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('legal_officer', 'admin', 'notary', 'nasabah'));

ALTER TABLE financing_applications
    ADD COLUMN IF NOT EXISTS user_id BIGINT NULL REFERENCES users (id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_financing_user ON financing_applications (user_id) WHERE deleted_at IS NULL;
