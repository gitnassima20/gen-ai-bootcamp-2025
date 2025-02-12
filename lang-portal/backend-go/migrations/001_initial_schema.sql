-- TODO: Add initial schema migration
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word TEXT NOT NULL,
    translation TEXT NOT NULL
    -- Add other columns as needed
);
