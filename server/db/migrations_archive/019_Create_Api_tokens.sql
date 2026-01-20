-- API Tokens for CLI Access
-- Phase 1: User Management / Phase 5: CLI Tool from roadmap
CREATE TABLE IF NOT EXISTS api_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    
    -- Token Information
    name TEXT NOT NULL, -- User-friendly name (e.g., "My Laptop", "CI/CD Pipeline")
    token_hash TEXT NOT NULL UNIQUE, -- Hashed token (never store plaintext)
    token_prefix TEXT NOT NULL, -- First 8 chars for identification (e.g., "mist_abc")
    
    -- Permissions & Scope
    scopes TEXT, -- JSON array of scopes: ["apps:read", "apps:write", "deployments:*", etc.]
    
    -- Usage Tracking
    last_used_at DATETIME,
    last_used_ip TEXT,
    usage_count INTEGER DEFAULT 0,
    
    -- Expiration
    expires_at DATETIME,
    
    -- Metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    revoked_at DATETIME, -- Soft revoke
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_api_tokens_user_id ON api_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_api_tokens_token_hash ON api_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_api_tokens_token_prefix ON api_tokens(token_prefix);
CREATE INDEX IF NOT EXISTS idx_api_tokens_expires_at ON api_tokens(expires_at);
