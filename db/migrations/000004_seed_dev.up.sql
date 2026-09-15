-- Migration 000004: seed minimum untuk development lokal.
-- Idempotent via ON CONFLICT DO NOTHING. Password admin: admin123 (bcrypt).

INSERT INTO users (full_name, email, password_hash, role, is_active, created_at, updated_at)
VALUES (
    'Administrator',
    'admin@bprs.local',
    '$2a$10$ITtSG0W2m9t8KWWREtGu8uUWHZ9xD.qrK8VpMXx4sf1Gu0yG1M.Vu',
    'admin',
    TRUE,
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO notaries (full_name, notary_number, wilayah_kerja, phone_number, email, is_available, created_at, updated_at)
VALUES (
    'Notaris Contoh, S.H.',
    'SK-001-DEV',
    'Jakarta Selatan',
    '081234567890',
    'notaris.contoh@bprs.local',
    TRUE,
    NOW(),
    NOW()
)
ON CONFLICT (notary_number) DO NOTHING;
