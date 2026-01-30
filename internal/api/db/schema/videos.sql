CREATE TABLE IF NOT EXISTS videos (
    id              TEXT PRIMARY KEY,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL,
    filename        TEXT NOT NULL,
    extension       TEXT NOT NULL,
    thumbnail       TEXT NOT NULL,
    duration        INTEGER NOT NULL,
    views           INTEGER NOT NULL,
    uploaded_at     TIMESTAMP NOT NULL
);
