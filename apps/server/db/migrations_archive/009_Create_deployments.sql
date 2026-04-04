CREATE TABLE IF NOT EXISTS deployments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_id INTEGER NOT NULL,
    
    -- Git Information
    commit_hash TEXT NOT NULL,
    commit_message TEXT,
    commit_author TEXT,
    
    -- Deployment Tracking
    triggered_by INTEGER,
    deployment_number INTEGER, -- Sequential number per app for easier reference
    
    -- Container Information
    container_id TEXT,
    container_name TEXT,
    image_tag TEXT, -- Store image tag for rollback support
    
    -- Build & Deploy Logs
    logs TEXT,
    build_logs_path TEXT, -- Path to build log file
    
    -- Status & Progress (Phase 1: Deployment improvements from roadmap)
    status TEXT NOT NULL CHECK(status IN ('pending', 'building', 'deploying', 'success', 'failed', 'stopped', 'rolled_back')) DEFAULT 'pending',
    stage TEXT DEFAULT 'pending', -- Current stage: cloning, building, deploying, etc.
    progress INTEGER DEFAULT 0, -- Percentage: 0-100
    error_message TEXT,
    
    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    started_at DATETIME,
    finished_at DATETIME,
    duration INTEGER, -- Duration in seconds
    
    -- Rollback Support (Phase 1: Critical feature from roadmap)
    is_active BOOLEAN DEFAULT 0, -- Currently running deployment
    rolled_back_from INTEGER, -- If this is a rollback, references the deployment it rolled back from
    
    FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
    FOREIGN KEY (triggered_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (rolled_back_from) REFERENCES deployments(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_deployments_app_id ON deployments(app_id);
CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments(status);
CREATE INDEX IF NOT EXISTS idx_deployments_is_active ON deployments(is_active);
CREATE INDEX IF NOT EXISTS idx_deployments_created_at ON deployments(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_deployments_commit_hash ON deployments(commit_hash);

