-- Migration 000008 down.
-- Catatan: rollback menghapus relasi user_id; pastikan tidak ada baris
-- financing_applications.user_id yang terisi sebelum rollback.
DROP INDEX IF EXISTS idx_financing_user;
ALTER TABLE financing_applications DROP COLUMN IF EXISTS user_id;

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('legal_officer', 'admin', 'notary'));
