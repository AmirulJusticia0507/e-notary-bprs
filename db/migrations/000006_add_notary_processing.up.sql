-- Migration 000006: Tahap 3 sistem notaris + Tahap 4 BPN (kolom processing di legal_documents).

ALTER TABLE legal_documents
    ADD COLUMN IF NOT EXISTS act_number TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS minutes_status TEXT NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS notary_fee BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS processing_status TEXT NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS notary_processed_at TIMESTAMPTZ NULL;
