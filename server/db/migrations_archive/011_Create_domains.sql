CREATE TABLE IF NOT EXISTS domains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_id INTEGER NOT NULL,
    domain_name TEXT NOT NULL UNIQUE,
    
    -- SSL/TLS Configuration (Phase 1: SSL/TLS Automation from roadmap)
    ssl_status TEXT NOT NULL CHECK(ssl_status IN ('pending', 'active', 'expired', 'failed', 'disabled')) DEFAULT 'pending',
    ssl_provider TEXT CHECK(ssl_provider IN ('letsencrypt', 'custom', 'none')) DEFAULT 'letsencrypt',
    
    -- Certificate Storage
    certificate_path TEXT,
    certificate_data TEXT, -- Store cert in DB as backup
    key_path TEXT,
    key_data TEXT, -- Store key in DB as backup (encrypted)
    chain_path TEXT, -- Certificate chain file
    
    -- Let's Encrypt Integration
    acme_account_url TEXT, -- ACME account URL
    acme_challenge_type TEXT CHECK(acme_challenge_type IN ('http-01', 'dns-01', 'tls-alpn-01')), -- Challenge type used
    
    -- Certificate Metadata
    issuer TEXT, -- Certificate issuer
    issued_at DATETIME, -- When certificate was issued
    expires_at DATETIME, -- Certificate expiry date
    last_renewal_attempt DATETIME,
    renewal_error TEXT,
    
    -- Configuration
    auto_renew BOOLEAN DEFAULT 1,
    force_https BOOLEAN DEFAULT 0, -- Redirect HTTP to HTTPS
    hsts_enabled BOOLEAN DEFAULT 0, -- HTTP Strict Transport Security
    hsts_max_age INTEGER DEFAULT 31536000, -- HSTS max-age in seconds (1 year default)
    
    -- WWW Redirect
    redirect_www BOOLEAN DEFAULT 0, -- Redirect www to non-www or vice versa
    redirect_www_to_root BOOLEAN DEFAULT 1, -- If true: www -> root, else: root -> www
    
    -- Metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_domains_app_id ON domains(app_id);
CREATE INDEX IF NOT EXISTS idx_domains_ssl_status ON domains(ssl_status);
CREATE INDEX IF NOT EXISTS idx_domains_expires_at ON domains(expires_at);
