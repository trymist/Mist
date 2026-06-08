CREATE TABLE IF NOT EXISTS environmentVariables (
    id TEXT PRIMARY KEY,
    -- resourceType can be 'project' or 'environment' or 'resource'
    resourceID TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(resourceID, key)
);

CREATE INDEX IF NOT EXISTS idx_environment_variables_resource_id ON environmentVariables(resourceID);
