-- User Sessions for Session Management
-- Phase 1: User Management from roadmap
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY, -- Session ID (UUID or secure random string)
    user_id INTEGER NOT NULL,
    
    -- Session Data
    session_data TEXT, -- JSON object containing session information
    
    -- Device & Location
    ip_address TEXT,
    user_agent TEXT,
    device_type TEXT CHECK(device_type IN ('desktop', 'mobile', 'tablet', 'unknown')),
    browser TEXT,
    os TEXT,
    location TEXT, -- City, Country based on IP (optional)
    
    -- Session Status
    is_active BOOLEAN DEFAULT 1,
    last_activity_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- Security
    revoked_at DATETIME, -- When session was revoked
    revoked_reason TEXT, -- Why session was revoked
    
    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_is_active ON sessions(is_active);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_last_activity ON sessions(last_activity_at);
