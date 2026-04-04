# System Requirements

Hardware and software requirements for running Mist.

## Minimum Requirements

### Small Deployments (1-5 applications)

- **CPU**: 1 vCPU
- **RAM**: 1GB minimum, 2GB recommended
- **Disk**: 10GB SSD
- **Network**: 100 Mbps
- **Cost**: ~$5-10/month VPS

::: tip Build Performance
The build process (compiling Go and bundling frontend) requires significant CPU. During installation, 1 vCPU may take 5-10 minutes to build.
:::

### Medium Deployments (5-20 applications)

- **CPU**: 2 vCPU
- **RAM**: 2-4GB
- **Disk**: 20-50GB SSD
- **Network**: 500 Mbps
- **Cost**: ~$12-24/month VPS

### Large Deployments (20+ applications)

- **CPU**: 4+ vCPU
- **RAM**: 8GB+
- **Disk**: 100GB+ SSD
- **Network**: 1 Gbps
- **Cost**: ~$40+/month VPS

## Software Requirements

### Operating System

- Ubuntu 20.04 LTS or later
- Debian 11 or later
- Any Linux distribution with:
  - systemd
  - Docker support
  - Modern kernel (5.x+)

### Docker

**Required** - Must be installed before running Mist installation script

- Docker Engine 20.10 or later
- Docker Compose v2 (for Traefik)
- Docker socket accessible at `/var/run/docker.sock`

::: warning Install Docker First
The Mist installation script does NOT install Docker. You must install Docker manually before running the script.

[Install Docker →](https://docs.docker.com/engine/install/)
:::

### Build Dependencies

Automatically installed by the installation script:
- Go 1.22.11
- Bun (JavaScript runtime)
- Build essentials (gcc, g++, make)
- Git, curl, wget, unzip

### Network

- **Required Ports**: 80, 443, 8080 must be available
- Port 80: HTTP entry point (Traefik)
- Port 443: HTTPS entry point (Traefik)
- Port 8080: Mist API and Dashboard
- Port 8081: Traefik Dashboard (optional)

- **Outbound internet access** for:
  - GitHub API (api.github.com)
  - Docker Hub (hub.docker.com)
  - Package downloads (Go, Bun, system packages)
  - Let's Encrypt ACME challenges (acme-v02.api.letsencrypt.org)

## Recommended Cloud Providers

### DigitalOcean

- **Droplet**: Basic ($12/month, 2GB RAM)
- Pre-installed Docker images available
- Good for most use cases

### Hetzner Cloud

- **CX21**: (€5/month, 2GB RAM)
- Excellent price/performance
- European data centers

### Linode

- **Nanode**: ($5/month, 1GB RAM)
- Good for small deployments
- Global data centers

### Vultr

- **Regular Performance**: ($6/month, 1GB RAM)
- Multiple locations worldwide

## Storage Considerations

### Installation Size

- **Mist binary**: 8-15MB
- **Frontend assets**: ~5MB
- **Go installation**: ~150MB (if not already installed)
- **Bun installation**: ~50MB (if not already installed)
- **Build dependencies**: ~100MB
- **Total installation**: ~500MB-1GB

### Database

- SQLite database grows slowly (~10MB per 1000 deployments)
- Initial database: ~100KB
- Recommend at least 1GB for database growth

### Git Repositories

- Cloned during deployment to `/var/lib/mist/repos/{app-id}/`
- Average repo: 50-500MB
- **Temporary** - Cleaned up after successful builds
- Multiple concurrent clones possible

### Docker Images

- Images stored by Docker Engine
- Size varies: 100MB - 2GB per image
- Each deployment creates new image tagged with commit hash
- Old images remain until manually pruned
- Recommendation: 
  - Small deployments: 10-20GB
  - Medium deployments: 50-100GB
  - Large deployments: 100GB+

### Logs

- Build logs stored in filesystem: `/var/lib/mist/logs/{deployment-id}/build.log`
- Average log file: 1-5MB
- Not automatically cleaned up
- Recommend 5-10GB for logs

### User Uploads

- Avatars stored in `/var/lib/mist/uploads/avatar/`
- Maximum size: 5MB per avatar
- Minimal storage impact

## Performance Benchmarks

### API Response Time

- Average: 20-50ms
- p95: < 100ms
- p99: < 200ms

### Concurrent Deployments

- Small (1 CPU): 1-2 concurrent
- Medium (2 CPU): 3-5 concurrent
- Large (4+ CPU): 10+ concurrent

### WebSocket Connections

- Hundreds of concurrent connections supported
- Minimal resource overhead

## Scaling Considerations

<div class="coming-soon-banner">
  <h4>🚧 Multi-Node Support Coming Soon</h4>
  <p>Horizontal scaling planned for Phase 4.</p>
</div>

### Current Limitations

- Single server deployment only
- Vertical scaling (upgrade server) required
- No load balancing across nodes

### Future Plans

- Docker Swarm support
- Kubernetes option
- Multi-server deployments
- Shared storage for builds
