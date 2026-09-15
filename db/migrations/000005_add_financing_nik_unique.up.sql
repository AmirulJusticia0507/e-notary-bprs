-- Migration 000005: unique NIK untuk sinkronisasi CBS/LOS.
-- Required by FinancingRepository.SyncFromCBS: ON CONFLICT (customer_nik).

CREATE UNIQUE INDEX IF NOT EXISTS uq_financing_applications_customer_nik
    ON financing_applications (customer_nik);
