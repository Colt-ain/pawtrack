-- Fix events table ID to use SERIAL for PostgreSQL auto-increment
-- Drop and recreate the sequence
DROP TABLE IF EXISTS events CASCADE;

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    dog_id INTEGER REFERENCES dogs(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL,
    note VARCHAR(255),
    at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_events_type ON events(type);
CREATE INDEX idx_events_at ON events(at);
CREATE INDEX idx_events_dog_id ON events(dog_id);
