CREATE TABLE IF NOT EXISTS app_configs (
    resourceId TEXT PRIMARY KEY NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    createdByUserId TEXT REFERENCES users(id) ON DELETE SET NULL,
    description TEXT,

    sourceType TEXT NOT NULL CHECK (sourceType IN ('github', 'gitlab', 'rawGit', 'dockerHub', 'customRegistry')),
    sourceId TEXT NOT NULL,

    gitBranch TEXT,
    gitRepoName TEXT,
    gitRepoOwner TEXT,
    gitCloneUrl TEXT,

    dockerImageName TEXT,
    dockerImageTag TEXT,

    autoDeploy INTEGER NOT NULL CHECK (autoDeploy IN (0, 1)) DEFAULT 0,
    port INTEGER,
    rootDir TEXT,

    dockerFilePath TEXT,
    cpuLimit REAL,
    memoryLimit INTEGER,
    restartPolicy TEXT NOT NULL CHECK (restartPolicy IN ('always', 'on-failure', 'no', 'unless-stopped')) DEFAULT 'always',

    healthCheckPath TEXT,
    healthCheckInterval INTEGER,
    healthCheckTimeout INTEGER,
    healthCheckRetries INTEGER,

    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_app_configs_created_by_user_id ON app_configs(createdByUserId);
CREATE INDEX IF NOT EXISTS idx_app_configs_gitCloneUrl ON app_configs(gitCloneUrl);
