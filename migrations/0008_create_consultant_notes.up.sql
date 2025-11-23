CREATE TABLE consultant_notes (
    id SERIAL PRIMARY KEY,
    consultant_id INTEGER NOT NULL REFERENCES users(id),
    dog_id INTEGER NOT NULL REFERENCES dogs(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_consultant_notes_consultant_id ON consultant_notes(consultant_id);
CREATE INDEX idx_consultant_notes_dog_id ON consultant_notes(dog_id);
CREATE INDEX idx_consultant_notes_created_at ON consultant_notes(created_at);
CREATE INDEX idx_consultant_notes_updated_at ON consultant_notes(updated_at);
