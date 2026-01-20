
INSERT INTO system_settings (key, value) VALUES ('wildcard_domain', '')
ON CONFLICT(key) DO NOTHING;

INSERT INTO system_settings (key, value) VALUES ('mist_app_name', 'mist')
ON CONFLICT(key) DO NOTHING;
