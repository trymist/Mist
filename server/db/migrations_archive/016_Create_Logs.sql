CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source TEXT NOT NULL CHECK(source IN ('app', 'system')),
    source_id INTEGER, -- app_id if source='app'
    message TEXT NOT NULL,
    level TEXT CHECK(level IN ('info', 'warn', 'error', 'debug')) DEFAULT 'info',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
