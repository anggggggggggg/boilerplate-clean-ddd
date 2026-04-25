-- Migration: 000001_create_deposits_table.down.sql
-- Rollback dari: 000001_create_deposits_table.up.sql
--
-- Jalankan dengan: make migrate-down
-- Urutan drop harus kebalikan dari urutan CREATE di file .up.sql
-- ============================================================================

-- Hapus trigger terlebih dahulu sebelum tabel
DROP TRIGGER  IF EXISTS trg_deposits_updated_at ON deposits;
DROP FUNCTION IF EXISTS set_updated_at();

-- Index akan otomatis terhapus saat tabel di-drop,
-- tapi explisit di sini untuk kejelasan
DROP INDEX IF EXISTS idx_deposits_user_status;
DROP INDEX IF EXISTS idx_deposits_status;
DROP INDEX IF EXISTS idx_deposits_user_id;

-- Hapus tabel utama
DROP TABLE IF EXISTS deposits;
