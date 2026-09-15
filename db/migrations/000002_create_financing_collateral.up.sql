-- Migration 000002: financing_applications & collaterals (sync CBS/LOS).

CREATE TABLE IF NOT EXISTS financing_applications (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    customer_name TEXT NOT NULL,
    customer_nik TEXT NOT NULL,
    financing_amount BIGINT NOT NULL DEFAULT 0 CHECK (financing_amount >= 0),
    collateral_type TEXT NOT NULL DEFAULT '' ,
    collateral_details TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_financing_status ON financing_applications (status) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS collaterals (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    financing_id BIGINT NOT NULL REFERENCES financing_applications (id) ON DELETE RESTRICT,
    type TEXT NOT NULL,
    details TEXT NOT NULL DEFAULT '',
    estimated_value BIGINT NOT NULL DEFAULT 0 CHECK (estimated_value >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_collaterals_financing ON collaterals (financing_id) WHERE deleted_at IS NULL;
