DROP TABLE IF EXISTS consultant_access;
DROP INDEX IF EXISTS idx_events_dog_id;
ALTER TABLE events DROP COLUMN dog_id;
