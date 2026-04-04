# Security Best Practices

Secure your Mist installation in production.

## Server Security

### Firewall Configuration

Configure firewall to allow only necessary ports:

```bash
# UFW example (Ubuntu/Debian)
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH (change 22 if using custom port)
sudo ufw allow 22/tcp

# Allow HTTP and HTTPS (required for web access and SSL challenges)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow Mist API (if accessing directly, otherwise only localhost)
sudo ufw allow 8080/tcp

# Optional: Allow Traefik dashboard (restrict to specific IPs in production)
sudo ufw allow 8081/tcp

# Enable firewall
sudo ufw enable

# Verify rules
sudo ufw status numbered
```

**Production Firewall Best Practices:**

```bash
# Restrict Mist API to localhost only (recommended if using Traefik)
sudo ufw delete allow 8080/tcp
sudo ufw allow from 127.0.0.1 to any port 8080

# Restrict Traefik dashboard to localhost only
sudo ufw delete allow 8081/tcp
sudo ufw allow from 127.0.0.1 to any port 8081

# Or allow only from specific admin IP
sudo ufw allow from YOUR_ADMIN_IP to any port 8081

# Allow SSH only from specific IP (very secure)
sudo ufw delete allow 22/tcp
sudo ufw allow from YOUR_ADMIN_IP to any port 22
```

### SSH Hardening

```bash
# Disable root login
sudo sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# Disable password auth (use SSH keys only)
sudo sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config

# Restart SSH
sudo systemctl restart sshd
```

### Keep System Updated

```bash
# Ubuntu/Debian
sudo apt update && sudo apt upgrade -y

# Enable automatic security updates
sudo apt install unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades
```

## Mist Security

### Strong Passwords

Enforce strong password requirements for all users:

- **Minimum length**: 12 characters (16+ recommended)
- **Complexity**: Mix of uppercase, lowercase, numbers, and symbols
- **No common passwords**: Avoid dictionary words, names, dates
- **Use password manager**: 1Password, Bitwarden, LastPass
- **Unique per account**: Never reuse passwords

**Owner Account Security:**
- The first user account becomes the Owner with full system access
- Use an extremely strong password for the Owner account
- Consider using a hardware security key (coming soon)

### JWT Secret

::: warning Hardcoded JWT Secret
Currently, the JWT secret is **hardcoded** in the application (`server/api/middleware/auth.go:14`). This is a **security risk** for production deployments.

**Workaround:** Modify the code before deployment:
```bash
cd /opt/mist/server/api/middleware
# Edit auth.go and change line 14 to use a strong random value
# Then rebuild: cd /opt/mist/server && go build -o ../bin/server main.go
```

**Upcoming:** Configurable JWT secret via environment variable.
:::

Generate a strong JWT secret:

```bash
# Generate random 64-character secret
openssl rand -hex 32

# Or use this one-liner
openssl rand -base64 48 | tr -d "=+/" | cut -c1-64
```

### Session Security

- **Session duration**: 31 days (hardcoded)
- **Sessions** are stored as JWT tokens
- **Token expiry**: Automatically logged out after 31 days
- **No session revocation**: Currently, there's no way to force-logout users (coming soon)

### GitHub Webhook Secret

GitHub webhook requests are validated using HMAC-SHA256 signatures.

**How it works:**
- Webhook secret is automatically generated when creating GitHub App
- Stored securely in Mist database
- All incoming webhook requests are validated before processing
- Invalid signatures are rejected

Generate a strong webhook secret if creating manually:

```bash
# Generate webhook secret (only needed for manual setup)
openssl rand -hex 32
```

Configure in: **Settings → Integrations → GitHub App** (automatic during app creation)

## Docker Security

### Socket Permissions

The Mist backend requires access to the Docker socket to manage containers:

```bash
# Verify Docker socket permissions
ls -la /var/run/docker.sock

# Should show: srw-rw---- 1 root docker

# Add your user to docker group (done during installation)
sudo usermod -aG docker $USER

# Log out and back in for group changes to take effect
```

::: danger Security Consideration
Access to the Docker socket = root-level access to the host system. The user running Mist can:
- Create privileged containers
- Mount host filesystem
- Escape container isolation

**Mitigations:**
- Run Mist as a non-root user (recommended, default in install script)
- Never expose Mist API publicly without authentication
- Use Docker rootless mode (advanced, requires configuration)
- Monitor Docker events for suspicious activity
:::

### Container Isolation

Applications deployed by Mist run in Docker containers with:

- **Network isolation**: Apps on `traefik-net` for controlled routing
- **Resource limits**: (coming soon - CPU/memory limits)
- **Read-only filesystems**: (coming soon)
- **Non-root users**: Depends on Dockerfile used

**Best Practices for Application Dockerfiles:**

```dockerfile
# Use specific versions, not :latest
FROM node:20-alpine

# Run as non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S appuser -u 1001
USER appuser

# Use multi-stage builds to reduce attack surface
# Copy only necessary files
# Scan images for vulnerabilities
```

### Image Security

- **Official base images**: Use official Docker images (node, python, golang, etc.)
- **Image scanning**: (coming soon - automatic vulnerability scanning)
- **Keep images updated**: Rebuild periodically with latest base images
- **Minimize layers**: Combine RUN commands to reduce image size
- **No secrets in images**: Never `COPY` or embed secrets in Dockerfiles

**Scan images manually:**

```bash
# Install Trivy (vulnerability scanner)
sudo apt install wget apt-transport-https gnupg lsb-release
wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | sudo apt-key add -
echo "deb https://aquasecurity.github.io/trivy-repo/deb $(lsb_release -sc) main" | sudo tee -a /etc/apt/sources.list.d/trivy.list
sudo apt update
sudo apt install trivy

# Scan an image
trivy image your-app-image:latest
```

## Network Security

### Reverse Proxy

- Always use Traefik as reverse proxy
- Never expose application ports directly
- Configure proper CORS headers

### SSL/TLS

Mist provides automatic SSL/TLS certificate provisioning and management.

**Current Features:**
- ✅ Automatic SSL certificate issuance via Let's Encrypt
- ✅ Automatic certificate renewal (90-day certificates, renewed at 60 days)
- ✅ HTTP to HTTPS redirect enabled by default for domains
- ✅ TLS 1.2+ encryption (configured by Traefik v3.1)
- ✅ Modern cipher suites
- ✅ Perfect Forward Secrecy (PFS)

**How it Works:**
1. User adds domain to application in Mist UI
2. Mist updates Traefik configuration with domain routing
3. Traefik automatically requests certificate from Let's Encrypt
4. Let's Encrypt verifies domain ownership via HTTP-01 challenge (port 80)
5. Certificate stored in `/opt/mist/letsencrypt/acme.json`
6. Traefik serves traffic over HTTPS with automatic HTTP redirect

**Certificate Storage:**
```bash
# View certificate file
sudo ls -la /opt/mist/letsencrypt/acme.json

# Should show: -rw------- (600 permissions)

# Backup certificates
sudo cp /opt/mist/letsencrypt/acme.json /var/backups/acme_$(date +%Y%m%d).json
```

**Verify SSL Configuration:**

```bash
# Test SSL connection
openssl s_client -connect your-domain.com:443 -servername your-domain.com

# Check certificate expiry
echo | openssl s_client -servername your-domain.com -connect your-domain.com:443 2>/dev/null | openssl x509 -noout -dates

# Use SSL Labs for comprehensive analysis
# https://www.ssllabs.com/ssltest/analyze.html?d=your-domain.com
```

[Learn more about SSL automation →](../guide/ssl-automation)

**Coming Soon:**
- Wildcard SSL certificates (`*.example.com`)
- Custom certificate upload
- HSTS headers and security policies
- Certificate expiry notifications
- mTLS (mutual TLS) for service-to-service communication

## Database Security

### SQLite File Permissions

The SQLite database contains ALL Mist data including secrets. Protect it carefully:

```bash
# Check current permissions
ls -la /var/lib/mist/mist.db

# Set restrictive permissions (read/write for owner only)
sudo chmod 600 /var/lib/mist/mist.db

# Ensure owned by the user running Mist (not root!)
sudo chown $USER:$USER /var/lib/mist/mist.db

# Verify
ls -la /var/lib/mist/mist.db
# Should show: -rw------- 1 youruser youruser
```

### Database Encryption

::: warning Database Not Encrypted
SQLite database is currently stored in **plaintext** on disk. This means:
- Anyone with file access can read the database
- All secrets (API keys, tokens, env vars) are visible
- Backups should be encrypted separately

**Mitigations:**
- Encrypt the entire disk/volume (LUKS, dm-crypt)
- Encrypt backups before uploading to cloud storage
- Restrict server access to trusted administrators only
- Use encrypted backup destinations (encrypted S3 buckets, etc.)
:::

**Encrypt backups:**

```bash
# Using GPG
gpg --symmetric --cipher-algo AES256 /var/backups/mist_backup.db

# Using OpenSSL
openssl enc -aes-256-cbc -salt -in /var/backups/mist_backup.db -out /var/backups/mist_backup.db.enc

# Decrypt when needed
openssl enc -d -aes-256-cbc -in /var/backups/mist_backup.db.enc -out /var/backups/mist_backup.db
```

### Backup Security

```bash
# Automated encrypted backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/var/backups/mist"
DB_FILE="/var/lib/mist/mist.db"
ENCRYPTION_KEY="your-encryption-passphrase"  # Store securely!

# Hot backup
sqlite3 $DB_FILE ".backup $BACKUP_DIR/mist_$DATE.db"

# Encrypt backup
openssl enc -aes-256-cbc -salt -pbkdf2 \
  -in $BACKUP_DIR/mist_$DATE.db \
  -out $BACKUP_DIR/mist_$DATE.db.enc \
  -pass pass:"$ENCRYPTION_KEY"

# Remove unencrypted backup
rm $BACKUP_DIR/mist_$DATE.db

# Keep last 30 days
find $BACKUP_DIR -name "mist_*.db.enc" -mtime +30 -delete
```

See [Backup & Recovery](./backup) for complete backup procedures.

## Application Security

### Environment Variables & Secrets

Environment variables can contain sensitive data (API keys, database passwords, etc.):

**Security Best Practices:**

```bash
# Never commit secrets to Git repositories
# Add .env to .gitignore

# Never log secrets
# Avoid console.log() or print() statements with env vars

# Rotate secrets regularly
# Change API keys every 90 days or after team member leaves

# Use strong random values
openssl rand -hex 32
```

**How Mist Handles Secrets:**

- Environment variables stored in SQLite database (plaintext)
- Passed to containers at deployment time via Docker `-e` flag
- Also passed as `--build-arg` during image builds
- Visible in `docker inspect` output on the host

::: warning Secrets Visibility
Environment variables are:
- **Stored in plaintext** in the database
- **Visible to anyone** with access to the Mist database or Docker host
- **Passed to build process** (avoid using secrets during build)

**Upcoming**: Secret encryption at rest and secret management integration (Vault, AWS Secrets Manager)
:::

**Minimize Secret Exposure:**

1. **Don't use secrets in Dockerfile ARG**: They're stored in image layers
   ```dockerfile
   # BAD - Secret visible in image history
   ARG API_KEY
   RUN curl -H "Authorization: $API_KEY" https://api.example.com/setup
   
   # GOOD - Fetch at runtime
   CMD node app.js  # App reads API_KEY from process.env at runtime
   ```

2. **Use secret management for databases**: Connect to external secret managers at runtime

3. **Limit access**: Only give secrets to apps that need them

### Audit Logging

Mist tracks important actions in the audit log:

**What's Logged:**
- User authentication (login/logout)
- Application creation/deletion
- Deployments
- Domain changes
- Environment variable changes
- Project member changes
- System settings changes

**Access Audit Logs:**
- Navigate to **Settings → Audit Logs** in the UI
- Or query database: `SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 100`

**Security Monitoring:**

```bash
# View recent audit logs via database
sqlite3 /var/lib/mist/mist.db "SELECT user_id, action, resource_type, created_at FROM audit_logs ORDER BY created_at DESC LIMIT 50;"
```

**What to Monitor:**
- Failed login attempts (possible brute force)
- Unusual deployment times (after hours activity)
- Unexpected user creations
- Permission changes
- Bulk deletions

::: tip
Set up alerts for critical actions (coming soon). For now, periodically review audit logs manually.
:::

## Monitoring & Incident Response

### Log Monitoring

Monitor system and application logs for security events:

```bash
# Watch Mist backend logs in real-time
sudo journalctl -u mist -f

# Filter for authentication events
sudo journalctl -u mist | grep -i "login\|auth\|failed"

# Watch Traefik access logs
docker logs -f traefik

# Monitor all Docker container activity
docker events

# Check for unauthorized container creation
docker ps -a --format "table {{.ID}}\t{{.Names}}\t{{.CreatedAt}}"
```

### System Monitoring

Monitor resource usage for anomalies:

```bash
# Monitor disk space (database growth, logs, containers)
df -h
du -sh /var/lib/mist/*
docker system df

# Monitor memory usage
free -h
docker stats --no-stream

# Monitor CPU usage
top -bn1 | head -20

# Check Docker daemon resource usage
ps aux | grep dockerd
```

### Security Alerts

Set up alerts for critical events (manual for now, automated coming soon):

**Monitor for:**
- Multiple failed login attempts
- Unusual deployment times
- High resource usage (potential crypto mining)
- New user creations
- Permission escalations
- Database file modifications
- Firewall rule changes

**Simple alerting script example:**

```bash
#!/bin/bash
# /usr/local/bin/mist-security-alert.sh

# Check for failed logins in last hour
FAILED_LOGINS=$(sudo journalctl -u mist --since "1 hour ago" | grep -c "authentication failed")

if [ $FAILED_LOGINS -gt 10 ]; then
    echo "WARNING: $FAILED_LOGINS failed login attempts in the last hour" | \
    mail -s "Mist Security Alert" admin@example.com
fi

# Check database size growth (possible data exfiltration)
DB_SIZE=$(stat -f%z /var/lib/mist/mist.db 2>/dev/null || stat -c%s /var/lib/mist/mist.db)
if [ $DB_SIZE -gt 1073741824 ]; then  # 1GB
    echo "WARNING: Database size exceeds 1GB" | \
    mail -s "Mist Database Alert" admin@example.com
fi
```

### Incident Response

If you suspect a security incident:

1. **Isolate**: Stop Mist and Traefik immediately
   ```bash
   sudo systemctl stop mist
   docker compose -f /opt/mist/traefik-compose.yml down
   ```

2. **Preserve Evidence**: Backup current state before changes
   ```bash
   sqlite3 /var/lib/mist/mist.db ".backup /tmp/incident_db_$(date +%Y%m%d_%H%M%S).db"
   sudo journalctl -u mist > /tmp/incident_logs_$(date +%Y%m%d_%H%M%S).log
   docker ps -a > /tmp/incident_containers_$(date +%Y%m%d_%H%M%S).txt
   ```

3. **Investigate**: Check audit logs, system logs, container logs

4. **Remediate**: Change passwords, rotate secrets, patch vulnerabilities

5. **Restore**: From clean backup if compromised

6. **Report**: If vulnerability in Mist itself, report to maintainers

## Coming Soon

The following security features are planned for future releases:

**Authentication & Authorization:**
- Two-Factor Authentication (TOTP/WebAuthn)
- SSO/SAML integration
- OAuth2 provider support
- API key management with scopes
- Session management (force logout, view active sessions)

**Secrets Management:**
- Environment variable encryption at rest
- Integration with HashiCorp Vault
- AWS Secrets Manager integration
- Secret rotation policies
- Secret scanning in code repositories

**Network Security:**
- Rate limiting per IP/user
- IP whitelisting for admin access
- Custom security headers (CSP, HSTS, X-Frame-Options)
- Web Application Firewall (WAF) integration
- DDoS protection configuration

**Application Security:**
- Automated vulnerability scanning for Docker images
- Security policy enforcement (required base images, etc.)
- Container resource limits (CPU, memory, network)
- Read-only container filesystems
- Seccomp profiles and AppArmor

**Monitoring & Compliance:**
- Security dashboard and metrics
- Automated security alerts (email, Slack, webhook)
- Compliance reporting (SOC 2, ISO 27001)
- Security audit trail export
- Real-time threat detection

**Data Protection:**
- Database encryption at rest
- Backup encryption by default
- Automated backup integrity verification
- Point-in-time recovery
- Data retention policies

## Security Checklist

Use this checklist to ensure your Mist installation is properly secured:

### Server Security
- [ ] Firewall configured and enabled (UFW/iptables)
- [ ] Only necessary ports open (22, 80, 443, 8080, 8081)
- [ ] Mist API (8080) and Traefik dashboard (8081) restricted to localhost or specific IPs
- [ ] SSH hardened (key-only auth, no root login)
- [ ] Automatic security updates enabled
- [ ] System packages up to date
- [ ] Disk encryption enabled (optional but recommended for sensitive data)

### Mist Security
- [ ] Strong passwords enforced for all users (12+ characters)
- [ ] Owner account uses extremely strong password
- [ ] JWT secret changed from hardcoded value (requires code modification)
- [ ] GitHub webhook secret configured and validated
- [ ] Regular security updates applied to Mist

### Docker Security
- [ ] Mist running as non-root user
- [ ] User in docker group (but aware of security implications)
- [ ] Docker socket permissions restricted (660)
- [ ] Containers use official base images only
- [ ] Images periodically scanned for vulnerabilities
- [ ] Unnecessary containers removed regularly

### Database Security
- [ ] Database file permissions locked down (600)
- [ ] Database owned by correct user (not root)
- [ ] Automated backups configured
- [ ] Backups stored securely (encrypted and offsite)
- [ ] Backup restoration tested regularly

### Network Security
- [ ] SSL/TLS enabled for all public domains
- [ ] Certificates automatically renewing
- [ ] HTTP automatically redirects to HTTPS
- [ ] Traefik dashboard access restricted
- [ ] No application ports exposed directly (all through Traefik)

### Application Security
- [ ] No secrets committed to Git repositories
- [ ] Environment variables contain only necessary data
- [ ] Secrets rotated regularly (90 days)
- [ ] Audit logs reviewed periodically
- [ ] Suspicious activity monitored

### Monitoring & Response
- [ ] Log monitoring configured
- [ ] Disk space alerts set up
- [ ] Resource usage monitored
- [ ] Incident response plan documented
- [ ] Emergency contacts updated

### Compliance (if applicable)
- [ ] Data retention policy defined
- [ ] User data handling documented
- [ ] Privacy policy published
- [ ] Terms of service published
- [ ] GDPR/compliance requirements met

## Reporting Security Issues

Found a security vulnerability in Mist?

**Please report responsibly:**

1. **DO NOT** open a public GitHub issue for security vulnerabilities
2. **DO** report via [GitHub Security Advisories](https://github.com/trymist/mist/security/advisories/new)
3. **OR** email: security@trymist.com (if available)
4. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

**We will:**
- Acknowledge receipt within 48 hours
- Provide estimated fix timeline
- Credit you in security advisory (if desired)
- Keep you updated on progress

Thank you for helping keep Mist secure!
