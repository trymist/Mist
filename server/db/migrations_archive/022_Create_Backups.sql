-- Database Backups
-- Phase 1: Backup & Recovery / Phase 2: Database Services from roadmap
CREATE TABLE IF NOT EXISTS backups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    -- Source
    app_id INTEGER NOT NULL,
    backup_type TEXT NOT NULL CHECK(backup_type IN ('manual', 'scheduled', 'pre_deployment', 'automatic')),
    
    -- Backup Details
    backup_name TEXT NOT NULL,
    file_path TEXT NOT NULL, -- Path to backup file
    file_size INTEGER, -- Size in bytes
    compression_type TEXT CHECK(compression_type IN ('none', 'gzip', 'bzip2', 'xz')) DEFAULT 'gzip',
    
    -- Database Information (for database service backups)
    database_type TEXT CHECK(database_type IN ('postgres', 'mysql', 'mariadb', 'mongodb', 'redis', 'other')),
    database_version TEXT,
    
    -- Storage Location
    storage_type TEXT CHECK(storage_type IN ('local', 's3', 'gcs', 'azure', 'ftp')) DEFAULT 'local',
    storage_path TEXT, -- Full path including bucket/container for cloud storage
    
    -- Backup Status
    status TEXT NOT NULL CHECK(status IN ('pending', 'in_progress', 'completed', 'failed', 'deleted')) DEFAULT 'pending',
    progress INTEGER DEFAULT 0, -- Percentage: 0-100
    error_message TEXT,
    
    -- Verification
    checksum TEXT, -- MD5 or SHA256 checksum
    checksum_algorithm TEXT CHECK(checksum_algorithm IN ('md5', 'sha256')) DEFAULT 'sha256',
    is_verified BOOLEAN DEFAULT 0,
    verified_at DATETIME,
    
    -- Restoration
    can_restore BOOLEAN DEFAULT 1,
    last_restore_at DATETIME,
    restore_count INTEGER DEFAULT 0,
    
    -- Retention
    retention_days INTEGER, -- NULL = keep forever
    auto_delete_at DATETIME, -- Calculated from retention_days
    
    -- Metadata
    created_by INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    duration INTEGER, -- Backup duration in seconds
    notes TEXT,
    
    FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_backups_app_id ON backups(app_id);
CREATE INDEX IF NOT EXISTS idx_backups_backup_type ON backups(backup_type);
CREATE INDEX IF NOT EXISTS idx_backups_status ON backups(status);
CREATE INDEX IF NOT EXISTS idx_backups_created_at ON backups(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_backups_auto_delete_at ON backups(auto_delete_at);
