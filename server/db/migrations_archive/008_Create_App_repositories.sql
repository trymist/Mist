CREATE TABLE IF NOT EXISTS app_repositories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_id INTEGER NOT NULL,
    source_type TEXT NOT NULL CHECK(source_type IN ('git_provider', 'github_app')),
    source_id INTEGER NOT NULL, -- references either git_providers.id OR github_installations.id
    repo_full_name TEXT NOT NULL,  -- e.g., vinayak/mist
    repo_url TEXT NOT NULL,
    branch TEXT DEFAULT 'main',
    webhook_id INTEGER,
    auto_deploy BOOLEAN DEFAULT 0,
    last_synced_at DATETIME,
    FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
    UNIQUE (app_id, repo_full_name)
);
