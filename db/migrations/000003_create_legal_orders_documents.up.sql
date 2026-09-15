-- Migration 000003: legal_orders, legal_order_logs, legal_documents (workflow SLA + audit trail).

CREATE TABLE IF NOT EXISTS legal_orders (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_number TEXT NOT NULL UNIQUE,
    financing_id BIGINT NOT NULL REFERENCES financing_applications (id) ON DELETE RESTRICT,
    notary_id BIGINT NOT NULL REFERENCES notaries (id) ON DELETE RESTRICT,
    assigned_to BIGINT NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'rejected')),
    sla_deadline TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_status ON legal_orders (status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_orders_notary ON legal_orders (notary_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_orders_assigned ON legal_orders (assigned_to) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_orders_financing ON legal_orders (financing_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS legal_order_logs (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES legal_orders (id) ON DELETE CASCADE,
    prev_status TEXT NOT NULL,
    new_status TEXT NOT NULL,
    changed_by BIGINT NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_order_logs_order ON legal_order_logs (order_id);

CREATE TABLE IF NOT EXISTS legal_documents (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES legal_orders (id) ON DELETE RESTRICT,
    file_url TEXT NOT NULL,
    sha256_hash CHAR(64) NOT NULL,
    e_meterai_sn TEXT NOT NULL DEFAULT '',
    e_sign_status TEXT NOT NULL DEFAULT 'pending' CHECK (e_sign_status IN ('pending', 'signed', 'failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_documents_order ON legal_documents (order_id) WHERE deleted_at IS NULL;
