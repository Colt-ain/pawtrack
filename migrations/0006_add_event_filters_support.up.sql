ALTER TABLE events ADD COLUMN dog_id INTEGER REFERENCES dogs(id);
CREATE INDEX idx_events_dog_id ON events(dog_id);

CREATE TABLE consultant_access (
    id SERIAL PRIMARY KEY,
    consultant_id INTEGER NOT NULL REFERENCES users(id),
    dog_id INTEGER NOT NULL REFERENCES dogs(id),
    granted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP
);
CREATE INDEX idx_consultant_access_consultant ON consultant_access(consultant_id);
CREATE INDEX idx_consultant_access_dog ON consultant_access(dog_id);
