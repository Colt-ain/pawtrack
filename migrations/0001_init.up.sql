-- initial schema

-- Postgres и SQLite совместимый синтаксис
CREATE TABLE IF NOT EXISTS events (
  id         INTEGER PRIMARY KEY,
  type       TEXT NOT NULL,
  note       TEXT,
  at         TIMESTAMP NOT NULL,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
CREATE INDEX IF NOT EXISTS idx_events_at   ON events(at);
