CREATE TABLE IF NOT EXISTS deployments (
		id TEXT PRIMARY KEY,
		resourceID TEXT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
		
		containerID TEXT,
		containerName TEXT,
		image TEXT,

		buildLogsPath TEXT,
		
		status TEXT NOT NULL CHECK (status IN (
			'pending',
			'building',
			'deploying',
			'success',
			'failed',
			'stopped',
			'rolledBack'
		)) DEFAULT 'pending',

		errorMessage TEXT,

		createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		startedAt DATETIME,
		completedAt DATETIME
		
);

CREATE INDEX IF NOT EXISTS idx_deployments_resource_id ON deployments(resourceID);
CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments(status);
