


CREATE TABLE IF NOT EXISTS projectMemberships (
    id TEXT PRIMARY KEY,
    projectId TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    userId TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
);

CREATE UNIQUE INDEX idx_project_memberships_project_user ON projectMemberships(projectId, userId);
