CREATE TABLE IF NOT EXISTS videos (
    id              TEXT PRIMARY KEY,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    extension       TEXT NOT NULL,
    duration        INTEGER NOT NULL DEFAULT 0,
    views           INTEGER NOT NULL DEFAULT 0,
    uploaded_at     TIMESTAMP NOT NULL
);
