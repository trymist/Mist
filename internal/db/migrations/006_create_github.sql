
CREATE TABLE IF NOT EXISTS githubApps (
    id TEXT PRIMARY KEY,
    githubAppName TEXT NOT NULL,
    githubAppId TEXT NOT NULL,
    githubClientId TEXT NOT NULL,
    githubClientSecret TEXT NOT NULL,
    githubPrivateKey TEXT, 
    githubWebhookSecret TEXT,
    createdByUserId TEXT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
);

CREATE INDEX idx_github_apps_created_by_user_id ON githubApps(createdByUserId);

CREATE TABLE IF NOT EXISTS githubInstallations (
    id TEXT PRIMARY KEY,
    githubAppId TEXT NOT NULL REFERENCES githubApps(id) ON DELETE CASCADE,
    githubInstallationId TEXT NOT NULL UNIQUE,
    accountName TEXT NOT NULL,
    accountType TEXT NOT NULL CHECK (accountType IN ('User', 'Organization')),
    userId TEXT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_github_installations_github_app_id ON githubInstallations(githubAppId);
CREATE INDEX idx_github_installations_user_id ON githubInstallations(userId);


