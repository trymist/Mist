CREATE TABLE IF NOT EXISTS git_providers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    provider TEXT NOT NULL CHECK(provider IN ('github', 'gitlab', 'bitbucket', 'gitea')),
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    expires_at DATETIME,
    username TEXT,
    email TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id, provider)
);
