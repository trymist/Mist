# Troubleshooting Guide

Common issues and solutions for Mist deployments.

::: tip Quick Links
- [Forgot Password Recovery](./forgot-password) - Reset your password using CLI
- [Manual Update Guide](./manual-update) - Update Mist manually if dashboard update fails
- [Community Support](https://discord.gg/kxK8XHR6) - Get help from the community
:::

---

## Installation Issues

### Installation Script Fails

**Problem:** Installation script exits with errors

**Solutions:**

1. **Check system requirements:**
```bash
# Verify Docker is installed and running
docker --version
docker ps

# Check disk space (need at least 2GB)
df -h /opt

# Verify network connectivity
curl -s https://github.com
```

2. **Run with proper permissions:**
```bash
sudo bash install.sh
```

3. **Check the installation log:**
```bash
tail -50 /tmp/mist-install.log
```

### Go Installation Fails

**Problem:** Go download or installation fails

**Solution:**
```bash
# Manually install Go 1.22+
wget https://go.dev/dl/go1.22.11.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.11.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

### Docker Permission Denied

**Problem:** Docker commands fail with permission errors

**Solution:**
```bash
# Add current user to docker group
sudo usermod -aG docker $USER

# Restart shell or run
newgrp docker

# Verify
docker ps
```

---

## Service Issues

### Mist Service Won't Start

**Problem:** `systemctl start mist` fails

**Diagnosis:**
```bash
# Check service status
sudo systemctl status mist

# View recent logs
sudo journalctl -u mist -n 50 --no-pager

# Check if port 8080 is already in use
sudo lsof -i :8080
```

**Solutions:**

1. **Port already in use:**
```bash
# Stop the conflicting service

# Reload and restart
sudo systemctl daemon-reload
sudo systemctl restart mist
```

2. **Database issues:**
```bash
# Check database file permissions
ls -la /var/lib/mist/mist.db

# Fix permissions if needed
sudo chown root:root /var/lib/mist/mist.db
sudo chmod 644 /var/lib/mist/mist.db
```

3. **Binary issues:**
```bash
# Rebuild the binary
cd /opt/mist/server
sudo go build -o mist
sudo systemctl restart mist
```

### Service Crashes After Update

**Problem:** Service won't start after running update

**Solution:**
```bash
# Check for stuck update logs
sudo journalctl -u mist -n 100 --no-pager | grep -i update

# If update is stuck, manually clear it using CLI
sudo mist-cli settings list

# View update logs via API or dashboard
curl http://localhost:8080/api/updates/history

# Force restart the service
sudo systemctl daemon-reload
sudo systemctl restart mist
```

---

## Deployment Issues

### Build Fails

**Problem:** Application build fails during deployment

**Diagnosis:**
```bash
# Check build logs in the dashboard
# Or check Docker logs
docker logs <container-id>
```

**Common causes:**

1. **Missing dependencies:**
```dockerfile
# Add missing dependencies to Dockerfile
RUN npm install
# or
RUN pip install -r requirements.txt
```

2. **Wrong build context:**
```dockerfile
# Ensure Dockerfile is in repository root
# Or specify custom path in deployment settings
```

3. **Build timeout:**
- Increase timeout in deployment settings
- Optimize Dockerfile with multi-stage builds
- Use build cache effectively

### Container Keeps Restarting

**Problem:** Deployed container restarts repeatedly

**Diagnosis:**
```bash
# View container logs
docker logs <container-name>

# Check container status
docker ps -a | grep <app-name>

# Inspect container
docker inspect <container-name>
```

**Common solutions:**

1. **Application crashes on startup:**
- Check environment variables are set correctly
- Verify database connection strings
- Check port configuration

2. **Health check failures:**
- Ensure health check endpoint returns 200
- Increase health check timeout
- Verify health check path is correct

3. **Missing environment variables:**
```bash
# Check via Mist dashboard or API
# Add missing variables in Environment Variables section
```

### Port Already in Use

**Problem:** Deployment fails with "port already allocated" error

**Solution:**
```bash
# Find conflicting container
docker ps | grep <port>

# Stop conflicting container
docker stop <container-id>

# Or change your app's port in environment variables
```

---

## Domain & SSL Issues

### Domain Not Resolving

**Problem:** Custom domain doesn't point to your application

**Checklist:**

1. **DNS Configuration:**
```bash
# Verify DNS records
dig yourdomain.com
nslookup yourdomain.com

# Should point to your server IP
# A record: yourdomain.com -> your-server-ip
```

2. **Wildcard domain setup:**
```bash
# Verify wildcard domain setting
mist-cli settings get --key wildcard_domain

# Should be set to your domain
mist-cli settings set --key wildcard_domain --value apps.example.com
```

3. **Traefik configuration:**
```bash
# Check Traefik is running
docker ps | grep traefik

# Restart Traefik if needed
cd /opt/mist
docker compose -f traefik-compose.yml restart
```

### SSL Certificate Issues

**Problem:** SSL certificate not generating or invalid

**Diagnosis:**
```bash
# Check Traefik logs
docker logs traefik

# Verify domain is accessible
curl -v https://yourdomain.com
```

**Solutions:**

1. **Let's Encrypt rate limit:**
- Wait 1 hour (for failures) or 1 week (for cert limit)
- Use staging environment for testing

2. **DNS not propagated:**
```bash
# Wait for DNS propagation (up to 48 hours)
# Check propagation status
dig +trace yourdomain.com
```

3. **Firewall blocking ports:**
```bash
# Ensure ports 80 and 443 are open
sudo ufw status
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

---

## Database Issues

### Cannot Connect to Database

**Problem:** Application can't connect to database service

**Solutions:**

1. **Verify database is running:**
```bash
docker ps | grep postgres
# or
docker ps | grep mysql
```

2. **Check connection string:**
- Ensure environment variables are correct
- Format: `postgresql://user:password@host:5432/dbname`
- Host should be the container name

3. **Network issues:**
```bash
# Verify containers are on same network
docker network inspect traefik-net
```

### Database Data Lost After Restart

**Problem:** Data disappears when container restarts

**Solution:**
- Ensure volumes are properly configured
- Check volume mounts:
```bash
docker inspect <container-name> | grep -A 10 Mounts
```

---

## GitHub Integration Issues

### GitHub App Not Connecting

**Problem:** Cannot connect GitHub repositories

**Diagnosis:**
1. Verify GitHub App is installed on your organization/account
2. Check webhook delivery in GitHub App settings
3. Ensure webhook URL is accessible from internet

**Solution:**
```bash
# Verify GitHub App settings in Mist dashboard
# Go to Settings -> Git

# Test webhook:
curl -X POST https://your-mist-domain.com/api/github/webhook \
  -H "Content-Type: application/json" \
  -d '{"action": "ping"}'
```

### Webhook Delivery Fails

**Problem:** GitHub webhooks not reaching Mist

**Solutions:**

1. **Firewall blocking webhooks:**
```bash
# Ensure Mist is accessible from internet
# Check if your server has public IP
curl ifconfig.me
```

2. **Wrong webhook URL:**
- Update in GitHub App settings
- Format: `https://your-mist-domain.com/api/github/webhook`

---

## Performance Issues

### Slow Dashboard Loading

**Problem:** Dashboard takes long to load

**Solutions:**

1. **Check system resources:**
```bash
# View CPU and memory usage
htop

# Check disk space
df -h

# View container resource usage
docker stats
```

2. **Too many logs:**
```bash
# Clean up old logs
sudo find /var/lib/mist/logs -type f -mtime +7 -delete
```

3. **Database optimization:**
```bash
# Vacuum SQLite database
sqlite3 /var/lib/mist/mist.db "VACUUM;"
```

### High CPU Usage

**Problem:** Server CPU constantly high

**Diagnosis:**
```bash
# Check which process is using CPU
top -bn1 | head -20

# Check Docker container CPU
docker stats --no-stream
```

**Solutions:**

1. **Too many containers:**
- Enable auto cleanup in settings
- Manually clean stopped containers:
```bash
docker container prune -f
```

2. **Container resource limits:**
- Set CPU/memory limits in deployment settings

---

## Update Issues

::: tip Manual Update Available
If you're experiencing issues with dashboard updates, you can [update manually using the install script](./manual-update).
:::

### Update Stuck in Progress

**Problem:** Update shows "in progress" but nothing happens

**Solution:**
```bash
# Check if update process is actually running
ps aux | grep install.sh

# If not running, the startup check should auto-resolve it
# Force restart to trigger the check
sudo systemctl restart mist

# Alternatively, manually clear stuck update via API
# (requires owner access)
curl -X POST http://localhost:8080/api/updates/clear?id=<update_log_id> \
  --cookie "session_token=your_token"
```

**Or perform a manual update:**
```bash
# See full manual update guide
curl -fsSL https://trymist.cloud/install.sh | sudo bash
```

[→ Full Manual Update Guide](./manual-update)

### Update Failed

**Problem:** Update failed and service won't start

**Solution:**
```bash
# Check update logs
sudo journalctl -u mist -n 100 --no-pager

# Revert to previous version if needed
cd /opt/mist
git log --oneline -10
sudo git reset --hard <previous-commit>
cd server
sudo go build -o mist
sudo systemctl restart mist
```

**Or perform a manual update with automatic recovery:**
```bash
curl -fsSL https://trymist.cloud/install.sh | sudo bash
```

[→ Full Manual Update Guide](./manual-update)

---

## CLI Issues

### CLI Command Not Found

**Problem:** `mist-cli` command not found

**Solution:**
```bash
# Check if CLI is installed
ls -la /usr/local/bin/mist-cli

# If not found, rebuild and install
cd /opt/mist/cli
sudo go build -o mist-cli
sudo cp mist-cli /usr/local/bin/
sudo chmod +x /usr/local/bin/mist-cli
```

### CLI Database Access Denied

**Problem:** CLI can't access database

**Solution:**
```bash
# Run with sudo
sudo mist-cli user list

# Or fix database permissions
sudo chown root:root /var/lib/mist/mist.db
sudo chmod 644 /var/lib/mist/mist.db
```

---

## Network Issues

### Cannot Access Dashboard

**Problem:** Cannot reach Mist dashboard in browser

**Checklist:**

1. **Service running:**
```bash
sudo systemctl status mist
curl http://localhost:8080/api/health
```

2. **Firewall:**
```bash
# Check firewall rules
sudo ufw status

# Allow port 8080
sudo ufw allow 8080/tcp
```

3. **Port forwarding:**
- If behind NAT, ensure port forwarding is configured
- Check router/cloud provider settings

### Applications Not Accessible

**Problem:** Deployed apps return 404 or 502 errors

**Diagnosis:**
```bash
# Check Traefik logs
docker logs traefik

# Verify Traefik config
cat /var/lib/mist/traefik/dynamic.yml

# Check app container is running
docker ps | grep <app-name>
```

**Solutions:**

1. **Traefik not running:**
```bash
cd /opt/mist
docker compose -f traefik-compose.yml up -d
```

2. **Wrong domain configuration:**
- Verify domain is correctly set in Mist dashboard
- Check DNS points to your server

3. **Container not in traefik-net:**
```bash
docker network connect traefik-net <container-name>
```

---

## Getting More Help

If you're still experiencing issues:

1. **Check logs:**
```bash
# Service logs
sudo journalctl -u mist -n 100 --no-pager

# Container logs
docker logs <container-name>

# Traefik logs
docker logs traefik
```

2. **Enable debug logging:**
```bash
# Edit service file
sudo nano /etc/systemd/system/mist.service

# Add debug environment variable
Environment=LOG_LEVEL=debug

# Restart
sudo systemctl daemon-reload
sudo systemctl restart mist
```

3. **Community support:**
- GitHub Issues: [github.com/trymist/mist/issues](https://github.com/trymist/mist/issues)
- Discord: [discord.gg/kxK8XHR6](https://discord.gg/hr6TCQDDkj)

4. **Gather diagnostic information:**
```bash
# System info
uname -a
docker version
docker compose version

# Mist version
curl http://localhost:8080/api/updates/version

# Resource usage
free -h
df -h
docker stats --no-stream
```
