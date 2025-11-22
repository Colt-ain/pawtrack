ALTER TABLE dogs ADD COLUMN owner_id INTEGER NOT NULL DEFAULT 1 REFERENCES users(id);
CREATE INDEX idx_dogs_owner_id ON dogs(owner_id);
