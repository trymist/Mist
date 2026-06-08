CREATE TABLE IF NOT EXISTS domains (
    id TEXT PRIMARY KEY,
    resourceID TEXT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    serviceName TEXT,
    domainName TEXT NOT NULL,
    port INTEGER NOT NULL,
    httpsEnabled BOOLEAN NOT NULL DEFAULT FALSE,
    certificateProvider TEXT NOT NULL CHECK (certificateProvider IN ('letsencrypt', 'custom')) DEFAULT 'letsencrypt', 
    certificatePath TEXT,
    customCertificateID TEXT,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_domains_resource_id ON domains(resourceID);
CREATE INDEX IF NOT EXISTS idx_domains_domain_name ON domains(domainName);
