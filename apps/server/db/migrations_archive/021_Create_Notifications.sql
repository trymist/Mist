-- Notifications System
-- Phase 1: Notification System from roadmap
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    -- Recipient
    user_id INTEGER, -- NULL for system-wide notifications
    
    -- Notification Type
    type TEXT NOT NULL CHECK(type IN (
        'deployment_success', 
        'deployment_failed', 
        'deployment_started',
        'ssl_expiry_warning',
        'ssl_renewal_success',
        'ssl_renewal_failed',
        'resource_alert',
        'app_error',
        'app_stopped',
        'backup_success',
        'backup_failed',
        'user_invited',
        'member_added',
        'system_update',
        'custom'
    )),
    
    -- Content
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    link TEXT, -- Link to relevant resource
    
    -- Related Resources
    resource_type TEXT, -- 'app', 'deployment', 'domain', 'project', etc.
    resource_id INTEGER,
    
    -- Channel-specific data
    email_sent BOOLEAN DEFAULT 0,
    email_sent_at DATETIME,
    
    slack_sent BOOLEAN DEFAULT 0,
    slack_sent_at DATETIME,
    
    discord_sent BOOLEAN DEFAULT 0,
    discord_sent_at DATETIME,
    
    webhook_sent BOOLEAN DEFAULT 0,
    webhook_sent_at DATETIME,
    
    -- Status
    is_read BOOLEAN DEFAULT 0,
    read_at DATETIME,
    
    priority TEXT CHECK(priority IN ('low', 'normal', 'high', 'urgent')) DEFAULT 'normal',
    
    -- Metadata
    metadata TEXT, -- JSON object with additional data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME, -- Optional expiry for temporary notifications
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_priority ON notifications(priority);
