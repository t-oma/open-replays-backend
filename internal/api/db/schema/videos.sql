CREATE TABLE IF NOT EXISTS videos (
    filename TEXT PRIMARY KEY,
    extension TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL
);
