-- Add DNS validation fields to domains table
ALTER TABLE domains ADD COLUMN dns_configured BOOLEAN DEFAULT 0;
ALTER TABLE domains ADD COLUMN dns_verified_at DATETIME;
ALTER TABLE domains ADD COLUMN last_dns_check DATETIME;
ALTER TABLE domains ADD COLUMN dns_check_error TEXT;

-- Add index for DNS status checks
CREATE INDEX IF NOT EXISTS idx_domains_dns_configured ON domains(dns_configured);
