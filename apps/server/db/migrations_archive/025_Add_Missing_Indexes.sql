
CREATE INDEX IF NOT EXISTS idx_envs_app_id ON envs(app_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_id ON audit_logs(resource_id);

CREATE INDEX IF NOT EXISTS idx_deployments_app_id ON deployments(app_id);

CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments(status);

CREATE INDEX IF NOT EXISTS idx_domains_app_id ON domains(app_id);

CREATE INDEX IF NOT EXISTS idx_project_members_user_id ON project_members(user_id);

CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members(project_id);

CREATE INDEX IF NOT EXISTS idx_apps_project_id ON apps(project_id);

CREATE INDEX IF NOT EXISTS idx_apps_status ON apps(status);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);

CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
