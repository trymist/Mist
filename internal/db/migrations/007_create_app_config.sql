

CREATE TABLE IF NOT EXISTS app_configs (
    resourceId TEXT PRIMARY KEY NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    createdByUserId TEXT REFERENCES users(id) ON DELETE SET NULL,
    description TEXT,

    sourceType TEXT NOT NULL CHECK (source IN ('github', 'gitlab', 'rawGit', 'dockerHub', 'customRegistry')),
    sourceId TEXT NOT NULL,

    -- git specific fields
    gitBranch TEXT,
    gitRepoName TEXT,
    gitRepoOwner TEXT,
    gitCloneUrl TEXT,

    -- docker specific fields
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
)
