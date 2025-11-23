-- Rollback: This would require restoring data, so just drop it
DROP TABLE IF EXISTS events CASCADE;

-- Recreate with old schema (not recommended, but for rollback)
CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    type TEXT NOT NULL,
    note TEXT,
    at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
CREATE INDEX IF NOT EXISTS idx_events_at ON events(at);
