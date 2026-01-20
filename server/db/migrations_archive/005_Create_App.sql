CREATE TABLE IF NOT EXISTS apps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    created_by INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    
    -- App Type (for our port handling solution)
    app_type TEXT NOT NULL CHECK(app_type IN ('web', 'service', 'database')) DEFAULT 'web',
    template_name TEXT, -- For pre-made services: 'postgres', 'redis', 'mongodb', etc.
    
    -- Git Configuration
    git_provider_id INTEGER,
    git_repository TEXT,
    git_branch TEXT DEFAULT 'main',
    git_clone_url TEXT,
    
    -- Deployment Configuration
    deployment_strategy TEXT NOT NULL CHECK(deployment_strategy IN ('auto', 'manual')) DEFAULT 'auto',
    port INTEGER,
    root_directory TEXT DEFAULT '.',
    build_command TEXT,
    start_command TEXT,
    dockerfile_path TEXT DEFAULT 'DOCKERFILE',
    
    -- Resource Limits (Phase 1: Resource Management from roadmap)
    cpu_limit REAL, -- CPU cores (e.g., 0.5, 1.0, 2.0)
    memory_limit INTEGER, -- Memory in MB (e.g., 512, 1024, 2048)
    restart_policy TEXT CHECK(restart_policy IN ('no', 'always', 'on-failure', 'unless-stopped')) DEFAULT 'unless-stopped',
    
    -- Health Check Configuration
    healthcheck_path TEXT,
    healthcheck_interval INTEGER DEFAULT 30,
    healthcheck_timeout INTEGER DEFAULT 10,
    healthcheck_retries INTEGER DEFAULT 3,
    
    -- Status
    status TEXT NOT NULL CHECK(status IN ('stopped', 'running', 'error', 'building', 'deploying')) DEFAULT 'stopped',
    
    -- Metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (git_provider_id) REFERENCES git_providers(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE (project_id, name)
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_apps_project_id ON apps(project_id);
CREATE INDEX IF NOT EXISTS idx_apps_app_type ON apps(app_type);
CREATE INDEX IF NOT EXISTS idx_apps_status ON apps(status);
