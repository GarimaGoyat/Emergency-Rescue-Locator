-- Location updates table: periodic GPS pings during active emergencies
CREATE TABLE IF NOT EXISTS location_updates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    emergency_id UUID NOT NULL REFERENCES emergencies(id) ON DELETE CASCADE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    accuracy DOUBLE PRECISION,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_location_updates_emergency_id ON location_updates(emergency_id);
CREATE INDEX IF NOT EXISTS idx_location_updates_recorded_at ON location_updates(recorded_at);
CREATE INDEX IF NOT EXISTS idx_location_updates_deleted_at ON location_updates(deleted_at);
