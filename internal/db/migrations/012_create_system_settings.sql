CREATE TABLE IF NOT EXISTS system_settings (
    key TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL,
    isSecret INTEGER NOT NULL DEFAULT 0,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);


INSERT INTO system_settings (key, value, isSecret) VALUES
('jwtSecret', 'mist-is-awesome', 1),
('version', '2.0.0', 0),
('appName', 'Mist', 0);

