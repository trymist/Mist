CREATE TABLE IF NOT EXISTS project_members (
   user_id INTEGER NOT NULL,
   project_id INTEGER NOT NULL,
   added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (user_id, project_id),
   FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
   FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
)
