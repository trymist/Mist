# Manual Update Guide

Step-by-step guide to manually update Mist when the dashboard update fails or gets stuck.

---

## Overview

If you're experiencing issues updating Mist through the dashboard, you can perform a manual update using the same installation script. This method is safe and will preserve all your data, applications, and settings.

::: tip Safe Update Process
The installation script is designed to be idempotent - you can run it multiple times safely. It will:
- ✅ Preserve your database and all data
- ✅ Keep all applications and deployments running
- ✅ Maintain your settings and configurations
- ✅ Update only the Mist binary and system files
:::

---

## When to Use Manual Update

Use manual update when:
- Dashboard update gets stuck in "in progress" state
- Update fails with errors in the dashboard
- Service won't restart after an update
- You want to update from a specific Git branch
- You need to rollback to a previous version

---

## Quick Update Steps

### 1. SSH into Your Server

```bash
ssh user@your-server-ip
```

### 2. Run the Installation Script

The same script used for installation will update your existing installation:

```bash
curl -fsSL https://trymist.cloud/install.sh | sudo bash
```

That's it! The script will:
1. Detect your existing installation
2. Update the repository to the latest version
3. Rebuild the binary
4. Restart the service
5. Verify everything is working

### 3. Verify Update

```bash
# Check service status
sudo systemctl status mist

# Check version
curl http://localhost:8080/api/updates/version

# Access dashboard
# https://your-mist-domain.com
```

---

## Detailed Manual Update Process

### Step 1: Prepare for Update

Before updating, it's good practice to check the current state:

```bash
# Check current version
curl http://localhost:8080/api/updates/version

# Check service status
sudo systemctl status mist

# View current running containers
docker ps
```

### Step 2: Run the Update Script

Execute the installation script with sudo:

```bash
curl -fsSL https://trymist.cloud/install.sh | sudo bash
```

**What the script does:**
1. Checks system requirements
2. Updates Git repository from the `release` branch
3. Downloads Go dependencies
4. Builds the new binary
5. Restarts the Mist service
6. Verifies the service is running

**Expected output:**
```
Starting Mist installation...
▶ Updating repository...
✔ Done
▶ Downloading dependencies...
✔ Done
▶ Building backend...
✔ Done
▶ Restarting service...
✔ Done
╔════════════════════════════════════════════╗
║ 🎉 Mist installation complete              ║
║ 👉 http://your-server-ip:8080              ║
╚════════════════════════════════════════════╝
```

### Step 3: Verify the Update

After the script completes:

1. **Check service is running:**
```bash
sudo systemctl status mist
```

Expected output: `Active: active (running)`

2. **Verify new version:**
```bash
curl http://localhost:8080/api/updates/version
```

3. **Check service logs:**
```bash
sudo journalctl -u mist -n 50 --no-pager
```

4. **Test dashboard access:**
   - Open your browser
   - Navigate to your Mist dashboard
   - Log in and verify functionality

5. **Check applications:**
   - Verify all applications are still running
   - Check deployment history
   - Test application access

---

## Handling Update Issues

### Update Script Fails

If the update script encounters errors:

**1. Check the installation log:**
```bash
tail -100 /tmp/mist-install.log
```

**2. Common issues:**

#### Port Already in Use
```bash
# Stop Mist service first
sudo systemctl stop mist

# Run update script
curl -fsSL https://trymist.cloud/install.sh | sudo bash
```

#### Git Conflicts
```bash
# Navigate to installation directory
cd /opt/mist

# Reset any local changes
sudo git reset --hard

# Run update script again
curl -fsSL https://trymist.cloud/install.sh | sudo bash
```

#### Build Failures
```bash
# Clear Go cache
sudo go clean -cache -modcache

# Rebuild manually
cd /opt/mist/server
sudo go mod tidy
sudo go build -o mist

# Restart service
sudo systemctl restart mist
```

### Service Won't Start After Update

If Mist service fails to start after update:

**1. Check service status and logs:**
```bash
sudo systemctl status mist
sudo journalctl -u mist -n 100 --no-pager
```

**2. Check for common issues:**

#### Database Locked
```bash
# Stop service
sudo systemctl stop mist

# Check for processes using database
sudo lsof /var/lib/mist/mist.db

# Kill any stale processes if found
sudo pkill -9 mist

# Start service
sudo systemctl start mist
```

#### Binary Not Executable
```bash
# Make binary executable
sudo chmod +x /opt/mist/server/mist

# Restart service
sudo systemctl restart mist
```

#### Missing Dependencies
```bash
# Reinstall dependencies
cd /opt/mist/server
sudo go mod download
sudo go build -o mist
sudo systemctl restart mist
```

### Applications Not Working After Update

If your applications are inaccessible after update:

**1. Check Traefik is running:**
```bash
docker ps | grep traefik
```

**2. Restart Traefik if needed:**
```bash
cd /opt/mist
docker compose -f traefik-compose.yml restart
```

**3. Check application containers:**
```bash
docker ps -a
```

**4. Restart specific containers if needed:**
```bash
docker restart <container-name>
```

---

## Update from Specific Version or Branch

### Update from Specific Git Branch

To update from a different branch (e.g., `main`, `develop`):

```bash
# Navigate to installation directory
cd /opt/mist

# Switch to desired branch
sudo git fetch origin
sudo git checkout <branch-name>
sudo git pull origin <branch-name>

# Build and restart
cd server
sudo go mod tidy
sudo go build -o mist
sudo systemctl restart mist
```

### Rollback to Previous Version

If the new update causes issues, you can rollback:

**1. View recent versions:**
```bash
cd /opt/mist
git log --oneline -10
```

**2. Rollback to specific commit:**
```bash
# Replace <commit-hash> with desired version
sudo git reset --hard <commit-hash>

# Rebuild
cd server
sudo go build -o mist

# Restart
sudo systemctl restart mist
```

**3. Rollback to previous release tag:**
```bash
# View available tags
git tag -l

# Checkout specific tag
sudo git checkout tags/v1.0.0

# Rebuild and restart
cd server
sudo go build -o mist
sudo systemctl restart mist
```

---

## Clearing Stuck Updates

If a dashboard update is stuck in "in progress" state:

### Method 1: Restart Service (Automatic Fix)

Simply restart the Mist service - the startup check will automatically resolve stuck updates:

```bash
sudo systemctl restart mist
```

The system will:
- Check for stuck updates on startup
- Compare current version with target version
- Mark update as success if versions match
- Mark as failed if versions don't match

### Method 2: Manual Fix Using CLI

Clear the stuck update manually:

**1. List update history:**
```bash
# Check update logs via dashboard or:
curl http://localhost:8080/api/updates/history
```

**2. Clear stuck update:**

Via API (requires authentication):
```bash
# Note: You need to be logged in and get the update_log_id
curl -X POST "http://localhost:8080/api/updates/clear?id=<update_log_id>" \
  --cookie "session_token=your_session_token"
```

### Method 3: Complete Manual Update

Perform a complete manual update which will automatically fix the stuck state:

```bash
# Run installation script
curl -fsSL https://trymist.cloud/install.sh | sudo bash

# Service restart will resolve stuck update
```

---

## Update Best Practices

### Before Updating

1. **Backup database:**
```bash
# Create backup directory
sudo mkdir -p /var/lib/mist/backups

# Copy database
sudo cp /var/lib/mist/mist.db /var/lib/mist/backups/mist.db.$(date +%Y%m%d_%H%M%S)
```

2. **Check disk space:**
```bash
df -h /opt
```

3. **Note current version:**
```bash
curl http://localhost:8080/api/updates/version
```

4. **Check running applications:**
```bash
docker ps
```

### During Update

1. **Monitor the update process:**
```bash
# In a separate terminal, watch logs
sudo journalctl -u mist -f
```

2. **Don't interrupt the process:**
   - Let the script complete
   - Don't cancel or close the terminal
   - Wait for confirmation message

### After Update

1. **Verify service health:**
```bash
# Check service
sudo systemctl status mist

# Check logs for errors
sudo journalctl -u mist -n 50 --no-pager | grep -i error

# Test health endpoint
curl http://localhost:8080/api/health
```

2. **Test critical functionality:**
   - Log into dashboard
   - Check applications are accessible
   - Test deployment (optional)
   - Verify domains are working

3. **Clean up old builds (optional):**
```bash
# Clean Docker resources if needed
docker system prune -f
```

---

## Automated Update Script

For convenience, you can create a custom update script:

```bash
#!/bin/bash
# save as: /usr/local/bin/mist-update

set -e

echo "🔄 Starting Mist update..."

# Backup database
echo "📦 Creating database backup..."
sudo mkdir -p /var/lib/mist/backups
sudo cp /var/lib/mist/mist.db /var/lib/mist/backups/mist.db.$(date +%Y%m%d_%H%M%S)

# Run update
echo "⬇️  Downloading and installing update..."
curl -fsSL https://trymist.cloud/install.sh | sudo bash

# Verify
echo "✅ Verifying update..."
sleep 3
sudo systemctl status mist --no-pager

echo "🎉 Update complete!"
echo "Current version:"
curl -s http://localhost:8080/api/updates/version | jq -r '.data.version'
```

**Make it executable and use it:**
```bash
sudo chmod +x /usr/local/bin/mist-update
sudo mist-update
```

---

## Update Frequency

### Recommended Update Schedule

- **Security updates:** As soon as available
- **Feature updates:** Monthly or as needed
- **Patch updates:** When issues are fixed

### Check for Available Updates

**Via Dashboard:**
- Navigate to Settings → Updates
- Click "Check for Updates"

**Via API:**
```bash
curl http://localhost:8080/api/updates/check
```

**Via GitHub:**
- Check [releases page](https://github.com/trymist/mist/releases)
- Watch repository for notifications

---

## Maintenance Mode (Optional)

For minimal downtime during updates:

**1. Create maintenance page:**
```bash
# Add to Traefik config to show maintenance page
# This is optional - most updates complete in under 30 seconds
```

**2. Notify users:**
- Post in team chat
- Send email notification
- Update status page

**3. Schedule updates:**
- Choose low-traffic periods
- Plan for 5-10 minutes downtime
- Have rollback plan ready

---

## Troubleshooting Specific Errors

### Error: "go: command not found"

```bash
# Install Go
wget https://go.dev/dl/go1.22.11.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.11.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Run update again
curl -fsSL https://trymist.cloud/install.sh | sudo bash
```

### Error: "Permission denied"

```bash
# Ensure you're using sudo
sudo bash install.sh

# Or fix ownership
sudo chown -R root:root /opt/mist
sudo chmod -R 755 /opt/mist
```

### Error: "Database is locked"

```bash
# Stop all processes
sudo systemctl stop mist

# Check for lock
sudo lsof /var/lib/mist/mist.db

# Kill if needed
sudo pkill -9 mist

# Restart
sudo systemctl start mist
```

---

## Getting Help

If manual update doesn't resolve your issue:

### Gather Information

```bash
# System info
uname -a

# Mist version
curl http://localhost:8080/api/updates/version

# Service status
sudo systemctl status mist

# Recent logs
sudo journalctl -u mist -n 100 --no-pager

# Installation log
tail -50 /tmp/mist-install.log
```

### Community Support

- **GitHub Issues:** [Report update issues](https://github.com/trymist/mist/issues)
- **Discord:** [Get real-time help](https://discord.gg/hr6TCQDDkj)
- **Discussions:** [Ask questions](https://github.com/trymist/mist/discussions)

### Include in Bug Report

When reporting update issues, include:
- Current version and target version
- Complete error messages
- Installation log (`/tmp/mist-install.log`)
- Service logs (`sudo journalctl -u mist -n 100`)
- System information (`uname -a`)

---

## Related Documentation

- [Upgrading Guide](../deployment/upgrading) - Official upgrade documentation
- [Common Issues](./) - General troubleshooting
- [Installation Guide](../deployment/installation) - Initial installation
- [Backup & Recovery](../deployment/backup) - Data backup procedures
