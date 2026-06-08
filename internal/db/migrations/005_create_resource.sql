CREATE TABLE IF NOT EXISTS resources (
    id TEXT PRIMARY KEY,
    environmentId TEXT NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('app', 'service')),
    description TEXT,
    status TEXT NOT NULL CHECK (status IN (
        'running',
        'stopped',
        'deploying',
        'failed',
        'deleted'
    )),
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_resources_environment_id ON resources(environmentId);

