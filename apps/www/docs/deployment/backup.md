# Backup & Recovery

Backup and restore your Mist installation.

## What to Backup

### 1. SQLite Database (Critical)

The database contains all your Mist data:
- Users, sessions, and authentication tokens
- Projects and team members
- Applications and configurations
- Deployments and deployment history
- Environment variables (including secrets)
- Domains and SSL certificate tracking
- GitHub App credentials
- Audit logs

**Location**: `/var/lib/mist/mist.db`

**Size**: Typically 10-100 MB depending on usage

::: danger Critical Data
This is the **single most important file** to backup. Without it, you lose all configuration and deployment history.
:::

### 2. SSL Certificates (Important)

Let's Encrypt SSL certificates for your domains.

**Location**: `/opt/mist/letsencrypt/acme.json`

**Note**: Certificates can be re-issued by Let's Encrypt, but backing them up prevents rate limiting and downtime during recovery.

### 3. Traefik Configuration (Auto-Generated)

Dynamic Traefik routing configuration.

**Location**: `/var/lib/mist/traefik/dynamic.yml`

**Note**: This file is auto-generated from the database, so restoring the database will recreate it.

### 4. User Uploads (Optional)

User avatars and uploaded files.

**Location**: `/var/lib/mist/uploads/`

**Size**: Usually small (< 100 MB)

### 5. Application Data (Important for Databases)

If you're running databases (PostgreSQL, MySQL, MongoDB, Redis) as applications in Mist, their data is stored in Docker volumes.

**List volumes:**
```bash
docker volume ls | grep mist
```

**Backup application volumes separately** if they contain important data.

### What NOT to Backup

- `/var/lib/mist/logs/` - Build logs (can be large, usually not needed)
- Git repositories - Cloned from source, can be re-cloned
- Docker images - Can be rebuilt
- `/var/lib/mist/traefik/dynamic.yml` - Auto-generated from database

## Manual Backup

### Quick Database Backup

```bash
# Stop Mist service for consistency
sudo systemctl stop mist

# Backup database with timestamp
sudo cp /var/lib/mist/mist.db /var/backups/mist_$(date +%Y%m%d_%H%M%S).db

# Start Mist service
sudo systemctl start mist

# Verify backup
ls -lh /var/backups/mist_*.db
```

::: tip Hot Backup (No Downtime)
SQLite supports online backups. You can backup without stopping Mist:

```bash
# Create backup directory
sudo mkdir -p /var/backups/mist

# Hot backup using SQLite
sqlite3 /var/lib/mist/mist.db ".backup /var/backups/mist_$(date +%Y%m%d_%H%M%S).db"
```
:::

### Full System Backup

Backup all Mist data including database, certificates, and uploads:

```bash
# Create backup directory
sudo mkdir -p /var/backups/mist

# Stop Mist for consistency (optional for hot backup)
sudo systemctl stop mist

# Backup everything
sudo tar -czf /var/backups/mist_full_$(date +%Y%m%d_%H%M%S).tar.gz \
  /var/lib/mist/mist.db \
  /var/lib/mist/uploads \
  /opt/mist/letsencrypt/acme.json \
  /opt/mist/traefik-static.yml

# Start Mist
sudo systemctl start mist

# Verify backup
ls -lh /var/backups/mist_full_*.tar.gz
```

### Backup Database Volumes (for Database Applications)

If you run databases (PostgreSQL, MySQL, etc.) as applications in Mist:

```bash
# List all Docker volumes
docker volume ls

# Backup a specific volume
docker run --rm -v app-123_data:/data -v /var/backups:/backup \
  alpine tar -czf /backup/app-123_data_$(date +%Y%m%d).tar.gz -C /data .
```

## Automated Backup Script

Create `/usr/local/bin/mist-backup.sh`:

```bash
#!/bin/bash

# Mist Automated Backup Script
BACKUP_DIR="/var/backups/mist"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30
DB_PATH="/var/lib/mist/mist.db"
CERT_PATH="/opt/mist/letsencrypt/acme.json"

# Create backup directory
mkdir -p $BACKUP_DIR

# Hot backup database using SQLite (no downtime)
if [ -f "$DB_PATH" ]; then
    sqlite3 $DB_PATH ".backup $BACKUP_DIR/mist_$DATE.db"
    echo "✅ Database backed up: mist_$DATE.db"
else
    echo "❌ Database not found at $DB_PATH"
    exit 1
fi

# Backup SSL certificates
if [ -f "$CERT_PATH" ]; then
    cp $CERT_PATH $BACKUP_DIR/acme_$DATE.json
    echo "✅ Certificates backed up: acme_$DATE.json"
fi

# Backup user uploads (if exists)
if [ -d "/var/lib/mist/uploads" ]; then
    tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz -C /var/lib/mist uploads
    echo "✅ Uploads backed up: uploads_$DATE.tar.gz"
fi

# Remove old backups (keep last N days)
find $BACKUP_DIR -name "mist_*.db" -mtime +$RETENTION_DAYS -delete
find $BACKUP_DIR -name "acme_*.json" -mtime +$RETENTION_DAYS -delete
find $BACKUP_DIR -name "uploads_*.tar.gz" -mtime +$RETENTION_DAYS -delete

echo "✅ Backup completed at $DATE"
echo "📁 Backup location: $BACKUP_DIR"

# Optional: Copy to remote server
# rsync -avz $BACKUP_DIR/ user@backup-server:/backups/mist/
```

Make executable and schedule:

```bash
# Make script executable
sudo chmod +x /usr/local/bin/mist-backup.sh

# Test the script
sudo /usr/local/bin/mist-backup.sh

# Schedule daily backups at 2 AM
(crontab -l 2>/dev/null; echo "0 2 * * * /usr/local/bin/mist-backup.sh") | crontab -

# Verify cron job
crontab -l
```

::: tip Backup Frequency
- **Development**: Daily backups with 7-day retention
- **Production**: Daily backups with 30-day retention + offsite copies
- **Critical systems**: Multiple daily backups + real-time replication
:::

## Backup to Cloud Storage

### AWS S3

```bash
# Install AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Configure credentials
aws configure
# Enter: Access Key ID, Secret Access Key, Region

# Backup to S3
aws s3 cp /var/backups/mist/mist_$(date +%Y%m%d_%H%M%S).db \
  s3://your-bucket/mist-backups/

# Automated S3 backup (add to backup script)
aws s3 sync /var/backups/mist/ s3://your-bucket/mist-backups/ --delete
```

### Backblaze B2

```bash
# Install B2 CLI
sudo pip3 install b2

# Authenticate
b2 authorize-account <keyId> <applicationKey>

# Create bucket (one-time)
b2 create-bucket mist-backups allPrivate

# Upload backup
b2 upload-file mist-backups /var/backups/mist/mist_latest.db mist_$(date +%Y%m%d).db
```

### Rsync to Remote Server

```bash
# Setup SSH key (one-time, run as backup user)
ssh-keygen -t ed25519 -f ~/.ssh/backup_key -N ""
ssh-copy-id -i ~/.ssh/backup_key user@backup-server

# Sync backups
rsync -avz -e "ssh -i ~/.ssh/backup_key" \
  /var/backups/mist/ \
  user@backup-server:/backups/mist/

# Add to automated script
rsync -avz --delete \
  /var/backups/mist/ \
  user@backup-server:/backups/mist/
```

### Rclone (Universal Cloud Backup)

Rclone supports 40+ cloud providers (S3, Google Drive, Dropbox, OneDrive, etc.):

```bash
# Install rclone
curl https://rclone.org/install.sh | sudo bash

# Configure cloud provider (interactive)
rclone config

# Backup to cloud
rclone copy /var/backups/mist/ remote:mist-backups/ -v

# Automated sync (add to backup script)
rclone sync /var/backups/mist/ remote:mist-backups/ --delete-excluded
```

## Restore from Backup

### Database Restore

```bash
# Stop Mist service
sudo systemctl stop mist

# Backup current database (just in case)
sudo cp /var/lib/mist/mist.db /var/lib/mist/mist.db.before-restore

# Restore database from backup
sudo cp /var/backups/mist/mist_20250122_140000.db /var/lib/mist/mist.db

# Fix permissions (important!)
sudo chown $USER:$USER /var/lib/mist/mist.db
sudo chmod 644 /var/lib/mist/mist.db

# Start Mist service
sudo systemctl start mist

# Verify service is running
sudo systemctl status mist

# Check logs for errors
journalctl -u mist -f
```

::: warning Permission Issues
Make sure the database file is owned by the user running the Mist service (usually the user who installed Mist, not `mist` or `root`).
:::

### Full System Restore

```bash
# Stop Mist service
sudo systemctl stop mist

# Extract full backup
sudo tar -xzf /var/backups/mist_full_20250122.tar.gz -C /

# Or restore individual files
sudo tar -xzf /var/backups/mist_full_20250122.tar.gz \
  -C / \
  var/lib/mist/mist.db \
  opt/mist/letsencrypt/acme.json

# Fix permissions
sudo chown -R $USER:$USER /var/lib/mist
sudo chmod 600 /opt/mist/letsencrypt/acme.json

# Restart services
sudo systemctl start mist
docker compose -f /opt/mist/traefik-compose.yml restart

# Verify everything is running
sudo systemctl status mist
docker ps | grep traefik
```

### Restore SSL Certificates

```bash
# Stop Traefik
docker compose -f /opt/mist/traefik-compose.yml down

# Restore certificates
sudo cp /var/backups/mist/acme_20250122.json /opt/mist/letsencrypt/acme.json
sudo chmod 600 /opt/mist/letsencrypt/acme.json

# Start Traefik
docker compose -f /opt/mist/traefik-compose.yml up -d

# Verify certificates are loaded
docker logs traefik 2>&1 | grep -i certificate
```

### Restore Docker Volumes (Database Applications)

```bash
# Stop the application container first
docker stop container-name

# Restore volume from backup
docker run --rm -v app-123_data:/data -v /var/backups:/backup \
  alpine sh -c "cd /data && tar -xzf /backup/app-123_data_20250122.tar.gz"

# Start the application container
docker start container-name
```

## Disaster Recovery

### Complete Server Loss

Follow these steps to recover from total server failure:

1. **Provision New Server**
   ```bash
   # Same OS as original (Ubuntu 20.04+ recommended)
   # Install Docker first
   curl -fsSL https://get.docker.com | sh
   ```

2. **Install Mist**
   ```bash
   curl -fsSL https://raw.githubusercontent.com/trymist/mist/main/install.sh | bash
   ```

3. **Stop Mist Service**
   ```bash
   sudo systemctl stop mist
   ```

4. **Restore Database**
   ```bash
   # Copy backup to new server
   scp backup-server:/backups/mist_latest.db /tmp/

   # Restore database
   sudo cp /tmp/mist_latest.db /var/lib/mist/mist.db
   sudo chown $USER:$USER /var/lib/mist/mist.db
   ```

5. **Restore SSL Certificates (optional, but recommended)**
   ```bash
   # Copy certificates
   scp backup-server:/backups/acme.json /tmp/

   # Restore
   sudo cp /tmp/acme.json /opt/mist/letsencrypt/acme.json
   sudo chmod 600 /opt/mist/letsencrypt/acme.json
   ```

6. **Start Services**
   ```bash
   sudo systemctl start mist
   docker compose -f /opt/mist/traefik-compose.yml restart
   ```

7. **Verify Recovery**
   ```bash
   # Check Mist service
   sudo systemctl status mist

   # Check Traefik
   docker ps | grep traefik

   # Access web interface
   curl http://localhost:8080/health  # or your domain
   ```

8. **Update DNS** (if IP changed)
   - Update A records to point to new server IP
   - Wait for DNS propagation (5-60 minutes)

### Database Corruption

Check and recover from database corruption:

```bash
# Check database integrity
sqlite3 /var/lib/mist/mist.db "PRAGMA integrity_check;"

# If output is "ok", database is fine
# If errors appear, restore from backup:

sudo systemctl stop mist
sudo cp /var/lib/mist/mist.db /var/lib/mist/mist.db.corrupted
sudo cp /var/backups/mist/mist_LATEST.db /var/lib/mist/mist.db
sudo systemctl start mist
```

### Accidental Data Deletion

If you accidentally delete data through the UI:

```bash
# Stop Mist immediately
sudo systemctl stop mist

# Restore from most recent backup
sudo cp /var/backups/mist/mist_$(ls -t /var/backups/mist/mist_*.db | head -1) \
  /var/lib/mist/mist.db

# Start Mist
sudo systemctl start mist
```

::: tip Recovery Point Objective (RPO)
Your RPO = time between backups. If you backup daily at 2 AM and disaster strikes at 1 PM, you lose up to 11 hours of data. For critical systems, increase backup frequency or use real-time replication.
:::

## Backup Best Practices

### The 3-2-1 Rule

- **3** copies of data
- **2** different storage types
- **1** offsite copy

### Backup Checklist

- [ ] Automated daily backups configured
- [ ] Backups stored on separate disk/server
- [ ] Offsite/cloud backup configured
- [ ] Backup retention policy defined (7-30 days)
- [ ] Backup integrity tested regularly
- [ ] Restore procedure documented
- [ ] Recovery time objective (RTO) defined
- [ ] Recovery point objective (RPO) defined

### Test Your Backups

Regularly test restore procedures:

```bash
# Test restore on separate server
scp backup-server:/backups/mist_latest.db /tmp/test-restore.db
sqlite3 /tmp/test-restore.db "SELECT COUNT(*) FROM users;"
```

## Coming Soon

The following backup features are planned for future releases:

- **Built-in Backup UI** - Configure and manage backups from the dashboard
- **One-Click Restore** - Restore from backup through the UI
- **Automated S3/Cloud Backup** - Built-in integration with cloud storage
- **Scheduled Backups** - Configure backup schedules per project or system-wide
- **Backup Verification** - Automatic integrity checks and restore testing
- **Point-in-Time Recovery** - Restore to a specific moment in time
- **Incremental Backups** - Save storage with differential backups
- **Database Replication** - Real-time database replication to standby server
- **Automated Application Volume Backups** - Backup database volumes automatically
- **Backup Encryption** - Encrypt backups at rest and in transit
- **Retention Policies** - Flexible retention rules (hourly, daily, weekly, monthly)
- **Webhook Notifications** - Get notified on backup success/failure

## Getting Help

If you encounter issues with backup or recovery:

- [GitHub Issues](https://github.com/trymist/mist/issues)
- [GitHub Discussions](https://github.com/trymist/mist/discussions)
