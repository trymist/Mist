CREATE TABLE IF NOT EXISTS update_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version_from TEXT NOT NULL,
    version_to TEXT NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('in_progress', 'success', 'failed')),
    logs TEXT,
    error_message TEXT,
    started_by INTEGER NOT NULL,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    FOREIGN KEY (started_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_update_logs_status ON update_logs(status);
CREATE INDEX IF NOT EXISTS idx_update_logs_started_at ON update_logs(started_at DESC);
