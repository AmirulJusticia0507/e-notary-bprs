-- Migration 000006 down.
ALTER TABLE legal_documents
    DROP COLUMN IF EXISTS notary_processed_at,
    DROP COLUMN IF EXISTS processing_status,
    DROP COLUMN IF EXISTS notary_fee,
    DROP COLUMN IF EXISTS minutes_status,
    DROP COLUMN IF EXISTS act_number;
