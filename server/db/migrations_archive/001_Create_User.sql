CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK(role IN ('owner', 'admin', 'user', 'viewer')) DEFAULT 'user',
    
    -- Profile
    full_name TEXT,
    avatar_url TEXT,
    bio TEXT,
    
    -- Authentication & Security (Phase 1: User Management from roadmap)
    email_verified BOOLEAN DEFAULT 0,
    email_verification_token TEXT,
    email_verification_sent_at DATETIME,
    
    password_reset_token TEXT,
    password_reset_expires_at DATETIME,
    password_changed_at DATETIME,
    
    two_factor_enabled BOOLEAN DEFAULT 0,
    two_factor_secret TEXT, -- TOTP secret
    two_factor_backup_codes TEXT, -- JSON array of backup codes
    
    -- Session & Login Tracking
    last_login_at DATETIME,
    last_login_ip TEXT,
    failed_login_attempts INTEGER DEFAULT 0,
    account_locked_until DATETIME,
    
    -- Preferences
    timezone TEXT DEFAULT 'UTC',
    language TEXT DEFAULT 'en',
    notification_preferences TEXT, -- JSON object of notification settings
    
    -- Metadata
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME -- Soft delete support
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
