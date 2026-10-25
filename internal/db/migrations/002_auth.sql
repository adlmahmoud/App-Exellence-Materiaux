CREATE TABLE IF NOT EXISTS auth_credentials (
    role TEXT PRIMARY KEY CHECK (role IN ('admin', 'user')),
    password_hash TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
