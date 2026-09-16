-- Migration 000010: foto profil user (nasabah saat signup).
-- Disimpan sebagai data-URL base64 (JPEG hasil downscale <=512px di frontend),
-- supaya tidak butuh storage file eksternal dan tetap awet di serverless.

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS photo_url TEXT NOT NULL DEFAULT '';
