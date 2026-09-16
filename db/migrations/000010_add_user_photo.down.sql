-- Migration 000010 down.
ALTER TABLE users DROP COLUMN IF EXISTS photo_url;
