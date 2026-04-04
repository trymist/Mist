# Architecture

Mist is designed as a lightweight, monolithic application with a focus on simplicity and resource efficiency.

## System Design

### High-Level Architecture

```
Internet Traffic
       ↓
   [Traefik Proxy] ← SSL/TLS, Domain Routing
       ↓
   [App Containers] ← Your deployed applications
       ↑
   [Mist Backend] ← Go API Server
       ↑
   [SQLite DB] ← Application state
       ↑
   [Frontend SPA] ← React dashboard
```

## Core Components

### 1. Backend API Server (Go)

**Location**: `/server`

The backend is a single Go binary that handles:

- **REST API** - HTTP endpoints for CRUD operations
- **WebSocket Server** - Real-time updates for logs and metrics
- **Authentication** - JWT-based auth with HTTP-only cookies
- **Deployment Queue** - In-memory queue for build jobs
- **Deployment Workers** - Goroutines that process deployments
- **Docker Integration** - Direct socket communication
- **GitHub Webhooks** - Auto-deploy on push events

**Key Features:**
- Single compiled binary
- No external process managers needed
- Embedded migrations
- Graceful shutdown handling

### 2. Frontend (React + Vite)

**Location**: `/dash`

Modern single-page application providing:

- Dashboard with real-time updates
- Application management interface
- Deployment monitoring
- Log viewer with WebSocket streaming
- System metrics visualization

**Tech Stack:**
- React 18
- Vite for build tooling
- Tailwind CSS for styling
- Shadcn/ui components

### 3. Database (SQLite)

**Location**: `/var/lib/mist/mist.db`

Embedded file-based database storing:

- Users and sessions
- Projects and members
- Applications
- Deployments
- Environment variables
- Domains
- Audit logs
- System settings

**Benefits:**
- No separate database server
- Automatic backups (just copy the file)
- ACID transactions
- Zero configuration

### 4. Docker Engine

Mist communicates with Docker via the socket (`/var/run/docker.sock`) to:

- Build images from Dockerfiles
- Create and manage containers
- Stream container logs
- Monitor container status
- Configure networks

### 5. Reverse Proxy (Traefik)

**Location**: `docker-compose.yml`

Traefik handles:

- HTTP/HTTPS traffic routing
- Domain-based routing to containers
- SSL/TLS termination
- Load balancing (future)

**Configuration:**
- Dynamic configuration via Docker labels
- File-based static configuration
- Automatic service discovery

## Deployment Workflow

### 1. User Triggers Deployment

```mermaid
User → Frontend → API → Deployment Queue
```

### 2. Queue Processing

```go
// Simplified deployment worker
func processDeployment(job DeploymentJob) {
    // 1. Clone Git repository
    cloneRepo(job.RepoURL, job.Branch)
    
    // 2. Build Docker image
    buildImage(job.Dockerfile, job.BuildArgs)
    
    // 3. Stop old container (if exists)
    stopContainer(job.AppID)
    
    // 4. Start new container
    startContainer(job.Image, job.EnvVars, job.Port)
    
    // 5. Update Traefik routes
    updateProxy(job.Domain, job.ContainerName)
    
    // 6. Update database
    markDeploymentComplete(job.ID)
}
```

### 3. Real-time Updates

WebSocket connections stream:
- Build logs during image creation
- Container logs after deployment
- System metrics every second
- Mist backend logs via journalctl integration

## Data Flow

### REST API Request

```
Frontend → API Endpoint → Middleware (Auth) → Handler → Database → Response
```

### WebSocket Connection

```
Frontend → WebSocket Upgrade → Auth Check → Stream Data → Client
```

### Deployment Flow

```
1. Git Push
2. GitHub Webhook → Mist API
3. Create Deployment Record
4. Add to Queue
5. Worker Picks Up Job
6. Clone Repo
7. Build Docker Image (stream logs via WebSocket)
8. Start Container
9. Update Traefik Config
10. Mark Deployment Complete
11. Notify Frontend via WebSocket
```

## Security Architecture

### Authentication

- JWT tokens stored in HTTP-only cookies
- Refresh token rotation
- Session management
- Password hashing with bcrypt

### Authorization

- Role-based access control (admin, user)
- Project-level permissions
- Owner-only operations

### API Security

- Rate limiting (planned)
- CORS configuration
- Input validation
- SQL injection prevention (parameterized queries)
- XSS protection

## Monitoring & Observability

### Logging Architecture

**Container Logs:**
- Streamed from Docker via socket connection
- Real-time delivery via WebSocket to frontend
- Last 100 lines buffered on connection

**System Logs:**
- Mist backend logs collected via systemd journal
- Accessed through `journalctl` integration
- Live streaming to `/logs` page
- Includes all Go backend output (API, deployments, errors)

**Build Logs:**
- Written to filesystem during Docker builds
- Path: `/var/lib/mist/logs/<deployment-id>.log`
- Accessible through deployment details

### Metrics Collection

**System Metrics:**
- CPU usage (percentage)
- Memory usage (used/total)
- Disk usage (used/total)
- Updated every second via WebSocket

**Container Metrics:**
- Container state (running/stopped/error)
- Container uptime
- Basic container information

::: tip
Advanced per-container metrics (CPU, memory, network I/O) are planned for future releases.
:::

## Scalability Considerations

### Current Design (Single Server)

✅ **Strengths:**
- Simple deployment
- No complex networking
- Easy to debug
- Low resource usage

⚠️ **Limitations:**
- Single point of failure
- Limited to one server's resources
- No horizontal scaling

### Future Improvements

🚧 **Planned:**
- Multi-node support with Docker Swarm
- Database replication
- Load balancing
- Shared storage for builds

## File System Layout

```
/var/lib/mist/
├── mist.db              # SQLite database
├── repos/               # Cloned Git repositories
│   └── <app-id>/       
│       └── <repo>/     
├── builds/              # Build artifacts
│   └── <deployment-id>/
└── logs/                # Build logs
    └── <deployment-id>.log

/etc/mist/
├── config.yml           # Mist configuration
└── traefik/             # Traefik configs
    ├── traefik.yml
    └── dynamic/

/usr/local/bin/
└── mist                 # Mist binary
```

## Performance Characteristics

### Resource Usage

**Idle State:**
- CPU: < 1%
- RAM: ~50-100MB
- Disk: ~100MB

**Under Load (10 concurrent deployments):**
- CPU: 10-30%
- RAM: ~200-500MB
- Disk: Depends on app sizes

### Benchmarks

- API Response Time: < 50ms (p95)
- WebSocket Latency: < 10ms
- Concurrent Users: 100+
- Apps per Instance: 1000+

## Technology Choices

### Why Go?

- Fast compilation and execution
- Excellent concurrency (goroutines)
- Single binary deployment
- Strong standard library
- Great Docker SDK

### Why SQLite?

- Zero configuration
- Embedded database
- ACID compliance
- Perfect for single-server deployments
- Easy backups (copy file)

### Why Traefik?

- Dynamic configuration
- Docker integration
- Automatic service discovery
- Let's Encrypt support
- Modern and actively maintained

## Learn More

- [Deployment Process](./deployments)
- [Traefik Configuration](../deployment/traefik)
