# Installation

Install Mist on your server with a single command.

## Prerequisites

### System Requirements

- **OS**: Ubuntu 20.04+, Debian 11+, Fedora, RHEL/CentOS, Arch Linux (any Linux with systemd)
- **RAM**: Minimum 1GB, recommended 2GB+
- **Disk**: Minimum 10GB free space
- **Docker**: Version 20.10 or later (**REQUIRED** - must be installed first)
- **Docker Compose**: v2 (usually bundled with Docker)
- **Ports**: 80, 443, 8080 must be available

### Docker Installation

::: warning Install Docker First
The Mist installation script does NOT install Docker. You must install Docker before running the script.
:::

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Log out and back in for group changes to take effect

# Verify installation
docker --version
docker compose version
```

[Official Docker Installation Guide →](https://docs.docker.com/engine/install/)

## Quick Install

### One-Line Installation

```bash
curl -fsSL https://trymist.cloud/install.sh | bash
```

During installation, you'll be prompted for:
- **Let's Encrypt Email**: Email for SSL certificate notifications (default: admin@example.com)

### What the Installation Script Does

The script performs the following steps:

1. **System Detection**
   - Detects package manager (apt/dnf/yum/pacman)
   - Checks for Docker installation

2. **Install Dependencies**
   - Git, curl, wget, unzip
   - Build tools (gcc, g++, make)

3. **Install Runtimes** (if not present)
   - Go 1.22.11 → `/usr/local/go`
   - Bun → `~/.bun`

4. **Clone Repository**
   - Clones Mist from GitHub → `/opt/mist`
   - Uses `main` branch

5. **Create Docker Network**
   - Creates `traefik-net` Docker bridge network

6. **Configure Traefik**
   - Updates `/opt/mist/traefik-static.yml` with your email
   - Starts Traefik via Docker Compose

7. **Build Frontend**
   - Runs `bun install` and `bun run build` in `/opt/mist/dash`
   - Copies built files to `/opt/mist/server/static/`

8. **Build Backend**
   - Runs `go mod tidy` and `go build` in `/opt/mist/server`
   - Creates binary at `/opt/mist/server/mist`

9. **Create Data Directory**
   - Creates `/var/lib/mist/`
   - Creates empty database file
   - Creates traefik config directory

10. **Configure Firewall**
    - Opens port 8080 (UFW/firewalld/iptables)

11. **Create Systemd Service**
    - Creates `/etc/systemd/system/mist.service`
    - Enables and starts service

::: tip Installation Time
On a 1 vCPU server, installation takes 5-10 minutes due to Go and frontend builds. On 2+ vCPUs, it takes 2-5 minutes.
:::

### Installation Paths

After installation, files are located at:

```
/opt/mist/                              # Main installation
├── server/
│   ├── mist                            # Go binary (executable)
│   └── static/                         # Frontend build files
├── dash/                               # Frontend source (kept)
├── traefik-static.yml                  # Traefik configuration
├── traefik-compose.yml                 # Docker Compose file
└── letsencrypt/                        # SSL certificates

/var/lib/mist/                          # Data directory
├── mist.db                             # SQLite database
├── traefik/                            # Dynamic configs
│   └── dynamic.yml                     # Generated at runtime
├── logs/                               # Build logs (created at runtime)
├── uploads/
│   └── avatar/                         # User avatars (created at runtime)
└── repos/                              # Git clones (temporary)

/etc/systemd/system/mist.service        # Service definition

/usr/local/go/                          # Go installation
~/.bun/                                 # Bun installation
```

## Post-Installation

### 1. Verify Installation

Check that Mist is running:

```bash
# Check service status
sudo systemctl status mist

# View logs
sudo journalctl -u mist -n 50 -f
```

### 2. Access Dashboard

Navigate to `http://your-server-ip:8080` in your browser.

::: tip First-Time Setup
On first visit, you'll create the owner account. This account has full system privileges.
:::

### 3. Create Admin Account

On the setup page:
1. Enter email address
2. Create a strong password
3. Click "Create Admin Account"
4. You'll be automatically logged in

### 4. Configure Wildcard Domain (Optional)

For automatic domain generation:
1. Go to **Settings** → **System**
2. Enter your wildcard domain (e.g., `example.com`)
3. Configure DNS with wildcard A record pointing `*.example.com` to your server
4. New web applications will automatically get domains like `{project}-{app}.example.com`

[Learn more about wildcard domains →](../guide/domains#wildcard-domain-configuration)

### 5. Configure GitHub Integration

For Git-based deployments:
1. Go to **Settings** → **Git**
2. Follow instructions to create a GitHub App
3. Install the app on your repositories

[Learn more about GitHub setup →](./github-app)

## Verification

### Check Components

Verify all components are running:

```bash
# Mist service
sudo systemctl status mist

# Traefik container
docker ps | grep traefik

# Docker network
docker network ls | grep traefik-net

# Database file
ls -lh /var/lib/mist/mist.db
```

### Test API

```bash
# Health check (should return system info)
curl http://localhost:8080/api/health
```

### Check Ports

Verify ports are accessible:

```bash
# Check if ports are listening
sudo netstat -tulpn | grep -E ':(80|443|8080) '

# Or with ss
sudo ss -tulpn | grep -E ':(80|443|8080) '
```

## Updating Installation

To update to the latest version:

```bash
# Re-run installation script
curl -sSL https://raw.githubusercontent.com/trymist/mist/main/install.sh | bash
```

::: warning Destructive Update
The installation script performs `git reset --hard origin/main`, which will overwrite any local changes. Always backup your data before updating.
:::

[Learn more about upgrading →](./upgrading)

## Uninstallation

To completely remove Mist:

```bash
curl -sSL https://raw.githubusercontent.com/trymist/mist/main/uninstall.sh | bash
```

This will:
- Stop and disable the systemd service
- Remove the service file
- Delete `/opt/mist/` directory
- Delete `/var/lib/mist/mist.db` database
- Remove firewall rules for port 8080
- Remove Go and Bun PATH entries from `~/.bashrc`

::: warning Data Preservation
The uninstall script does NOT remove:
- Traefik container (still running)
- Docker network `traefik-net`
- SSL certificates in `/opt/mist/letsencrypt/`
- `/var/lib/mist/traefik/` directory
- Go and Bun installations

Manually remove these if desired.
:::

## Troubleshooting

### Installation Fails

**Docker not installed:**
```
Error: Docker is not installed
```
Install Docker first using the official installation guide.

**Port already in use:**
```
Error: Port 8080 is already in use
```
Check what's using the port:
```bash
sudo lsof -i :8080
# or
sudo netstat -tulpn | grep :8080
```

**Build fails:**
```
Error: Go build failed
```
Check system resources and try again. Build requires significant CPU and RAM.

### Service Won't Start

Check logs for errors:

```bash
sudo journalctl -u mist -n 100 --no-pager
```

Common issues:
- Database file permissions
- Port 8080 already in use
- Docker socket not accessible
- Missing dependencies

### Docker Socket Permission Denied

If Mist can't access Docker:

```bash
# Add user to docker group
sudo usermod -aG docker $(whoami)

# Log out and back in, or:
newgrp docker

# Restart service
sudo systemctl restart mist
```

### Traefik Not Running

Check Traefik status:

```bash
docker ps | grep traefik
docker logs traefik
```

Restart Traefik:

```bash
cd /opt/mist
docker compose -f traefik-compose.yml down
docker compose -f traefik-compose.yml up -d
```

### Can't Access Dashboard

1. **Check firewall**:
   ```bash
   sudo ufw status
   sudo ufw allow 8080/tcp
   ```

2. **Check service**:
   ```bash
   sudo systemctl status mist
   ```

3. **Check if port is listening**:
   ```bash
   curl http://localhost:8080
   ```

4. **Check from external**:
   ```bash
   curl http://YOUR_SERVER_IP:8080
   ```

### Database Migration Errors

If database migrations fail:

```bash
# Backup database
sudo cp /var/lib/mist/mist.db /var/lib/mist/mist.db.backup

# Delete and recreate (CAUTION: loses all data)
sudo rm /var/lib/mist/mist.db
sudo touch /var/lib/mist/mist.db
sudo chown $(whoami):$(whoami) /var/lib/mist/mist.db

# Restart service
sudo systemctl restart mist
```

## Advanced Installation

### Custom Installation Path

The installation script uses `/opt/mist` by default. To use a custom path, clone and build manually:

```bash
# Clone to custom location
git clone https://github.com/trymist/mist.git /your/custom/path

# Follow manual build steps
cd /your/custom/path
# ... (frontend and backend build)

# Update systemd service WorkingDirectory
sudo nano /etc/systemd/system/mist.service
```

### Development Installation

For local development:

```bash
# Clone repository
git clone https://github.com/trymist/mist.git
cd mist

# Frontend
cd dash
bun install
bun run dev  # Runs on port 5173

# Backend (separate terminal)
cd server
go mod tidy
go run main.go  # Runs on port 8080
```

[Learn more about development setup →](https://github.com/trymist/mist#development)

## Security Recommendations

After installation:

1. **Change default firewall rules** - Only allow necessary ports
2. **Set strong password** - Use a password manager
3. **Configure GitHub webhook secret** - Secure webhook endpoint
4. **Regular backups** - Backup `/var/lib/mist/mist.db`
5. **Update regularly** - Keep Mist up to date
6. **Disable Traefik dashboard** - In production, disable port 8081
7. **Monitor logs** - Watch for suspicious activity

[Learn more about security →](./security)

## Next Steps

- [Configure GitHub Integration](./github-app)
- [Set Up Custom Domains](../guide/domains)
- [Create Your First Application](../guide/applications)
- [Configure Backups](./backup)
- [Review Security Best Practices](./security)
