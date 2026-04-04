# Upgrading Mist

Keep your Mist installation up to date with the latest features and security patches.

## Before You Upgrade

::: danger Always Backup First
**NEVER upgrade without a recent backup.** Database migrations can fail or cause data loss.

```bash
# Quick backup before upgrading
sudo systemctl stop mist
sqlite3 /var/lib/mist/mist.db ".backup /var/backups/mist_pre_upgrade_$(date +%Y%m%d_%H%M%S).db"
sudo systemctl start mist
```

See [Backup & Recovery](./backup) for complete backup procedures.
:::

### Pre-Upgrade Checklist

- [ ] Backup database (`/var/lib/mist/mist.db`)
- [ ] Backup SSL certificates (`/opt/mist/letsencrypt/acme.json`)
- [ ] Check [CHANGELOG.md](https://github.com/trymist/mist/blob/main/CHANGELOG.md) for breaking changes
- [ ] Note current deployed applications and their status
- [ ] Have rollback plan ready
- [ ] Schedule upgrade during low-traffic window

## Check Current Version

Check your current Mist version:

**Via Dashboard:**
1. Navigate to **Extras → Updates**
2. View current version at the top
3. Click "Check for Updates" to see if updates are available

**Via API:**
```bash
curl http://localhost:8080/api/updates/version
```

**Via CLI:**
```bash
# Check git commit (if available)
cd /opt/mist && git log -1 --oneline

# Check systemd service status
sudo systemctl status mist
```

## Upgrade Methods

### Method 1: Dashboard One-Click Update (Recommended)

Mist includes a built-in update system accessible from the dashboard:

1. Navigate to **Extras → Updates** in the dashboard
2. Click **"Check for Updates"** to see if a new version is available
3. Review the release notes and changes
4. Click **"Update Now"** to trigger the automatic update
5. The system will:
   - Download and install the latest version
   - Restart the Mist service
   - Run database migrations automatically
   - Mark the update as complete

::: tip Auto-Recovery
If an update appears stuck (rare), simply restart the Mist service:
```bash
sudo systemctl restart mist
```

The system automatically detects incomplete updates on startup and resolves them based on the installed version.
:::

::: warning Manual Clearing
Owners can manually clear stuck updates via the dashboard or API if needed:
```bash
curl -X POST http://localhost:8080/api/updates/clear \
  -H "Authorization: Bearer YOUR_TOKEN"
```
:::

### Method 2: Re-run Installation Script

The installation script is idempotent and will update your existing installation:

```bash
# 1. Backup database
sudo systemctl stop mist
sqlite3 /var/lib/mist/mist.db ".backup /var/backups/mist_backup_$(date +%Y%m%d).db"
sudo systemctl start mist

# 2. Re-run installation script
curl -fsSL https://raw.githubusercontent.com/trymist/mist/main/install.sh | bash

# 3. Verify upgrade
sudo systemctl status mist
sudo journalctl -u mist -n 50
```

**What it does:**
- Pulls latest code from GitHub (`main` branch)
- Rebuilds Go backend
- Rebuilds Vite frontend
- Restarts Mist service
- Preserves database and all data
- Runs database migrations automatically

### Method 3: Manual Git Pull and Rebuild

If you want more control over the upgrade process:

```bash
# 1. Backup database
sudo systemctl stop mist
sqlite3 /var/lib/mist/mist.db ".backup /var/backups/mist_backup_$(date +%Y%m%d).db"

# 2. Pull latest code
cd /opt/mist
git fetch origin
git pull origin main

# 3. Rebuild backend
cd server
go build -o ../bin/server main.go

# 4. Rebuild frontend
cd ../dash
npm install
npm run build

# 5. Restart service
sudo systemctl start mist

# 6. Verify
sudo systemctl status mist
sudo journalctl -u mist -n 50
```

### Method 4: Upgrade to Specific Version

To upgrade to a specific release or branch:

```bash
# 1. Backup
sudo systemctl stop mist
sqlite3 /var/lib/mist/mist.db ".backup /var/backups/mist_backup_$(date +%Y%m%d).db"

# 2. Checkout specific version
cd /opt/mist
git fetch --all --tags
git checkout tags/v1.2.3  # or specific commit hash

# 3. Rebuild (same as Method 2)
cd server && go build -o ../bin/server main.go
cd ../dash && npm install && npm run build

# 4. Restart
sudo systemctl start mist
```

## Database Migrations

Mist automatically runs database migrations on startup. The migration system is located in `server/db/migrations/`.

### How Migrations Work

1. **On Startup**: Mist checks which migrations have been applied
2. **Auto-Execute**: Runs any pending migrations in order
3. **Track Applied**: Stores migration history in the database
4. **Continue**: Starts normally after migrations complete

### Monitoring Migrations

Check migration status in logs:

```bash
# View recent Mist logs
sudo journalctl -u mist -n 100

# Follow logs in real-time during upgrade
sudo journalctl -u mist -f

# Search for migration-related logs
sudo journalctl -u mist | grep -i migration
sudo journalctl -u mist | grep -i "Running migration"
```

Look for messages like:
- `Running migrations...`
- `Migration 001_Create_User.sql completed`
- `All migrations applied successfully`

### Migration Locations

Migrations are stored in:
```
/opt/mist/server/db/migrations/
├── 001_Create_User.sql
├── 002_Create_Projects.sql
├── 003_Create_Project_members.sql
├── 004_Create_GitProviders.sql
└── ...
```

### If Migrations Fail

If a migration fails during upgrade:

```bash
# Stop Mist
sudo systemctl stop mist

# Check logs for error details
sudo journalctl -u mist -n 100 | grep -i error

# Restore from backup
sudo cp /var/backups/mist_backup_YYYYMMDD.db /var/lib/mist/mist.db

# Try again or report issue on GitHub
sudo systemctl start mist
```

::: warning
Never manually modify the database schema. Always use the built-in migration system or report schema issues on GitHub.
:::

## Post-Upgrade Verification

After upgrading, verify everything works correctly:

### 1. Check Service Status

```bash
# Verify Mist is running
sudo systemctl status mist

# Check for errors in recent logs
sudo journalctl -u mist -n 50
```

### 2. Test Web Interface and Version

```bash
# Test localhost access
curl http://localhost:8080/health

# Check updated version
curl http://localhost:8080/api/updates/version

# Or browse to your domain
# https://your-mist-domain.com
```

**Via Dashboard:**
- Log in and check the sidebar shows the new version
- Navigate to **Extras → Updates** to verify the update was successful

### 3. Verify Applications

- Log in to Mist dashboard
- Check all applications are showing as "Running"
- Test one or two applications to ensure they're accessible
- Check recent deployments in audit logs

### 4. Check Traefik

```bash
# Verify Traefik is running
docker ps | grep traefik

# Check Traefik logs for errors
docker logs traefik 2>&1 | tail -50
```

### 5. Test SSL Certificates

```bash
# Check certificate status for your domains
curl -I https://your-app-domain.com

# Verify SSL in browser (no warnings)
```

## Rollback Procedure

If the upgrade causes issues and you need to rollback:

### Quick Rollback

```bash
# 1. Stop Mist
sudo systemctl stop mist

# 2. Restore previous code (if available)
cd /opt/mist
git log --oneline -10  # Find previous commit
git reset --hard <previous-commit-hash>

# 3. Restore database backup
sudo cp /var/backups/mist_backup_YYYYMMDD.db /var/lib/mist/mist.db

# 4. Rebuild (if code changed)
cd server && go build -o ../bin/server main.go

# 5. Restart
sudo systemctl start mist

# 6. Verify
sudo systemctl status mist
```

### Full Rollback

If you need to completely restore to pre-upgrade state:

```bash
# Stop services
sudo systemctl stop mist
docker compose -f /opt/mist/traefik-compose.yml down

# Restore database
sudo cp /var/backups/mist_backup_YYYYMMDD.db /var/lib/mist/mist.db

# Restore certificates (if needed)
sudo cp /var/backups/acme_backup_YYYYMMDD.json /opt/mist/letsencrypt/acme.json

# Restore code to previous version
cd /opt/mist
git reset --hard <previous-commit>

# Rebuild
cd server && go build -o ../bin/server main.go
cd ../dash && npm install && npm run build

# Restart services
sudo systemctl start mist
docker compose -f /opt/mist/traefik-compose.yml up -d

# Verify
sudo systemctl status mist
```

::: danger Report Issues
If you need to rollback, please report the issue on [GitHub Issues](https://github.com/trymist/mist/issues) with:
- Error logs from `journalctl -u mist`
- What version you upgraded from/to
- Steps that led to the problem
:::

## Upgrade Best Practices

### Scheduling Upgrades

- **Development**: Upgrade anytime
- **Staging**: Upgrade and test thoroughly before production
- **Production**: Schedule during maintenance windows or low-traffic periods

### Testing Strategy

1. **Test in Staging First**: If possible, test upgrade on staging server
2. **Read Changelog**: Check for breaking changes and new features
3. **Backup Everything**: Database, certificates, uploads
4. **Gradual Rollout**: Test one application before full deployment
5. **Monitor Closely**: Watch logs and metrics for 24-48 hours after upgrade

### Stay Updated

- Watch [GitHub Releases](https://github.com/trymist/mist/releases)
- Subscribe to [GitHub Discussions](https://github.com/trymist/mist/discussions)
- Check for security updates regularly

### Recommended Upgrade Frequency

- **Security patches**: Apply immediately
- **Minor versions**: Monthly or when new features needed
- **Major versions**: Test thoroughly in staging first

## Troubleshooting Upgrades

### Common Issues

**Build fails during upgrade:**
```bash
# Clear build cache and retry
cd /opt/mist/server
go clean -cache
go build -o ../bin/server main.go

cd /opt/mist/dash
rm -rf node_modules package-lock.json
npm install
npm run build
```

**Service won't start after upgrade:**
```bash
# Check detailed logs
sudo journalctl -u mist -n 100 --no-pager

# Check for port conflicts
sudo netstat -tulpn | grep 8080

# Verify file permissions
ls -la /opt/mist/bin/server
```

**Database migration stuck:**
```bash
# Check if database is locked
sudo lsof /var/lib/mist/mist.db

# If locked, stop all processes accessing it
sudo systemctl stop mist
# Wait a few seconds
sudo systemctl start mist
```

**Update appears stuck in dashboard:**
```bash
# The system auto-recovers on restart
sudo systemctl restart mist

# Check logs to verify resolution
sudo journalctl -u mist -n 50

# Or manually clear via API
curl -X POST http://localhost:8080/api/updates/clear \
  -H "Authorization: Bearer YOUR_TOKEN"
```

See the [Troubleshooting Guide](../troubleshooting/) for more update-related issues.

**Applications not accessible after upgrade:**
```bash
# Restart Traefik
docker compose -f /opt/mist/traefik-compose.yml restart

# Check Traefik logs
docker logs traefik 2>&1 | tail -100

# Verify dynamic config
cat /var/lib/mist/traefik/dynamic.yml
```

## Available Now

The following upgrade features are currently available:

- ✅ **One-Click Dashboard Updates** - Update directly from **Extras → Updates**
- ✅ **Version Display** - View current version in sidebar and updates page
- ✅ **Update Status Tracking** - Monitor update progress and history
- ✅ **Automatic Recovery** - System auto-resolves stuck updates on restart
- ✅ **Manual Update Clear** - Owners can manually clear stuck updates
- ✅ **Check for Updates** - View available updates from dashboard
- ✅ **Database Auto-Migrations** - Migrations run automatically on update

## Coming Soon

The following upgrade features are planned for future releases:

- **In-App Update Notifications** - Push notifications when updates are available
- **Automatic Scheduled Updates** - Opt-in automatic updates with configurable schedule
- **Update Channels** - Choose stable, beta, or nightly update channel
- **Pre-Upgrade Health Check** - Automatic system checks before upgrade
- **One-Click Rollback** - Built-in rollback from dashboard
- **Staged Rollouts** - Upgrade one server at a time in multi-server setup
