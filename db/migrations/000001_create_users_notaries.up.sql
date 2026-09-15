-- Migration 000001: users & notaries (master data + soft delete)
-- Konvensi: PK BIGINT IDENTITY, created_at/updated_at TIMESTAMPTZ, deleted_at NULL, FK ON DELETE RESTRICT.

CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    full_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'legal_officer' CHECK (role IN ('legal_officer', 'admin', 'notary')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS notaries (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    full_name TEXT NOT NULL,
    notary_number TEXT NOT NULL UNIQUE,
    wilayah_kerja TEXT NOT NULL,
    phone_number TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL UNIQUE,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_notaries_available ON notaries (is_available) WHERE deleted_at IS NULL;
