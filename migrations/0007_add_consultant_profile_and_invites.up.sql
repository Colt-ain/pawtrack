CREATE TABLE consultant_profiles (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    description TEXT,
    services TEXT,
    breeds TEXT,
    location VARCHAR(255),
    surname VARCHAR(255)
);

CREATE TABLE invites (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users(id),
    consultant_id INTEGER NOT NULL REFERENCES users(id),
    dog_id INTEGER NOT NULL REFERENCES dogs(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_invites_token ON invites(token);
CREATE INDEX idx_invites_consultant_id ON invites(consultant_id);
