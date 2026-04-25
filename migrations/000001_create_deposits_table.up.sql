-- Migration: 000001_create_deposits_table.up.sql
-- Dibuat dengan: make migrate-create name=create_deposits_table
--
-- Tabel deposits adalah representasi persistence dari domain Deposit.
-- UUID digunakan sebagai primary key untuk menghindari sequential ID yang
-- dapat dieksploitasi (OWASP: Insecure Direct Object Reference).
-- ============================================================================

-- Aktifkan ekstensi uuid-ossp jika belum ada (untuk gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS deposits (
    id          UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID            NOT NULL,
    amount      NUMERIC(18, 2)  NOT NULL CHECK (amount > 0),

    status      VARCHAR(20)     NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'success', 'failed')),
    note        TEXT,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- ── Index ────────────────────────────────────────────────────────────────────
-- Query paling umum: cari semua deposit milik user tertentu
CREATE INDEX idx_deposits_user_id ON deposits (user_id);

-- Query monitoring: filter berdasarkan status (pending, success, failed)
CREATE INDEX idx_deposits_status ON deposits (status);

-- Query kombinasi: deposit user dengan status tertentu (misal: pending per user)
CREATE INDEX idx_deposits_user_status ON deposits (user_id, status);

-- ── Trigger: auto-update updated_at ─────────────────────────────────────────
-- Karena GORM `autoUpdateTime` bekerja di application layer,
-- tambahkan trigger sebagai fallback jika ada update langsung via SQL.
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_deposits_updated_at
    BEFORE UPDATE ON deposits
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
