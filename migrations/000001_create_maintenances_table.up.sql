-- =============================================================================
-- Migration: Create maintenances table
-- SAD Reference: ADR-2 — PostgreSQL via Supabase
-- Pattern: Database Per Service
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS maintenances (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id   UUID        NOT NULL,
    incident_id  UUID,                          -- NULL for preventive maintenance
    type         VARCHAR(20) NOT NULL CHECK (type IN ('corrective', 'preventive')),
    severity     SMALLINT    NOT NULL DEFAULT 0 CHECK (severity >= 0 AND severity <= 10),
    status       VARCHAR(20) NOT NULL DEFAULT 'queued'
                              CHECK (status IN ('queued', 'in_progress', 'completed', 'failed')),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- Index for queue queries (Process Network 3)
CREATE INDEX IF NOT EXISTS idx_maintenances_status ON maintenances (status);

-- Index for vehicle-based lookups
CREATE INDEX IF NOT EXISTS idx_maintenances_vehicle_id ON maintenances (vehicle_id);

-- Index for creation date ordering
CREATE INDEX IF NOT EXISTS idx_maintenances_created_at ON maintenances (created_at DESC);
