-- Rollback: Drop maintenances table
DROP INDEX IF EXISTS idx_maintenances_created_at;
DROP INDEX IF EXISTS idx_maintenances_vehicle_id;
DROP INDEX IF EXISTS idx_maintenances_status;
DROP TABLE IF EXISTS maintenances;
