# Applications

Applications are the core deployable units in Mist. Each application represents a containerized service that can be a web application, background service, or managed database.

## Application Types

Mist supports four distinct application types:

### Web Applications (`web`)

Web applications are services exposed via HTTP/HTTPS through Traefik reverse proxy:

- **Public access** via custom domains
- **Auto-generated domains** (if wildcard DNS configured)
- **SSL certificates** managed automatically
- **Port configuration** (default: 3000)
- **Git-based deployments** from GitHub repositories

**Use cases**: Websites, web apps, REST APIs, GraphQL servers

### Service Applications (`service`)

Background services that run without external HTTP access:

- **Internal only** - Not exposed to public internet
- **Port 3000** (internal, not routable)
- **Git-based deployments** from GitHub repositories
- **Same deployment flow** as web apps

**Use cases**: Workers, queue processors, scheduled tasks, internal services

### Database Applications (`database`)

Managed database services using pre-configured Docker templates:

- **Template-based** deployment (PostgreSQL, MySQL, Redis, MongoDB)
- **No Git repository** required
- **Pre-configured** CPU and memory limits from template
- **Version control** via Docker image tags
- **Environment variables** for configuration (passwords, databases, etc.)

**Use cases**: PostgreSQL, MySQL, Redis, MongoDB, or any other database service

### Docker Compose Applications (`compose`)

Multi-container applications defined by a `docker-compose.yml` file:

- **Multi-service** - Deploy multiple containers as a single app
- **Docker Compose** - Uses native `docker compose up` command
- **Git-based deployments** - Clones repo containing docker-compose.yml
- **Aggregated status** - Shows running count for all services
- **Combined logs** - View logs from all containers

**Use cases**: Complex stacks, microservices, apps with databases/caches, development environments

[Learn more about Docker Compose apps →](./docker-compose)

## Creating an Application

### Via Dashboard

1. Navigate to your project
2. Click **"New Application"**
3. Select application type:
   - **Web** - For public-facing applications
   - **Service** - For background workers
   - **Database** - For managed database services
   - **Docker Compose** - For multi-container applications
4. Fill in the required fields based on type

### Web/Service Configuration

**Initial Creation (Required):**
- **Name**: Unique identifier within project (lowercase, numbers, hyphens only)
- **Description**: Optional application description
- **Port**: Port your app listens on (default: 3000)

**After Creation (Configure in App Settings):**
- **Git Repository**: Connect GitHub repository
- **Git Branch**: Select branch to deploy from
- **Dockerfile Path**: Path to your Dockerfile (default: `./Dockerfile`)
- **Root Directory**: Root directory for builds (default: `/`)
- **Environment Variables**: Add key-value pairs in Environment tab

::: tip Build & Start Commands Not Yet Available
Build and start commands are coming soon. Currently, you must provide a Dockerfile for your application.
:::

### Database Configuration

- **Name**: Database instance name
- **Description**: Optional description
- **Template**: Select from available templates:
  - PostgreSQL
  - MySQL
  - Redis
  - MongoDB
- **Environment Variables**: Database credentials and settings

::: tip Auto-Generated Domains
For web applications, Mist can automatically generate a subdomain if wildcard DNS is configured in system settings. The format is `{project-name}-{app-name}.{wildcard-domain}`.

Example: If wildcard domain is `example.com`, project is `production`, and app is `api`, the auto-generated domain will be `production-api.example.com`.

This only applies to web applications. Service and database apps do not get auto-generated domains.
:::

## Resource Configuration

### CPU Limits

Set CPU core limits for your application:

- **Default**: No limit (uses available CPU)
- **Range**: 0.1 to N cores (e.g., 0.5, 1.0, 2.0)
- **Database defaults**: Applied automatically from templates

```javascript
// API example
{
  "appId": 1,
  "cpuLimit": 1.0  // 1 CPU core
}
```

### Memory Limits

Set memory limits in megabytes:

- **Default**: No limit (uses available memory)
- **Range**: 128MB to system maximum
- **Database defaults**: Applied automatically from templates

```javascript
// API example
{
  "appId": 1,
  "memoryLimit": 512  // 512 MB
}
```

### Restart Policies

Control container restart behavior:

- **`no`**: Never restart automatically
- **`always`**: Always restart if stopped
- **`on-failure`**: Restart only on failure
- **`unless-stopped`**: Restart unless manually stopped (default)

```javascript
// API example
{
  "appId": 1,
  "restartPolicy": "unless-stopped"
}
```

## Container Configuration

### Dockerfile Requirements

**All web and service applications require a Dockerfile.** Mist builds your application using Docker.

#### Basic Dockerfile Structure

Your repository must include a Dockerfile with:
1. **Base image** (FROM statement)
2. **Working directory** (WORKDIR)
3. **Dependency installation** (COPY package files, RUN install commands)
4. **Application code** (COPY source files)
5. **Port exposure** (EXPOSE - optional but recommended)
6. **Start command** (CMD or ENTRYPOINT)

Example for Node.js:
```dockerfile
FROM node:18-alpine
WORKDIR /app

# Copy dependency files
COPY package*.json ./

# Install dependencies
RUN npm ci --production

# Copy application code
COPY . .

# Expose port (should match port in Mist config)
EXPOSE 3000

# Start command
CMD ["node", "server.js"]
```

Example for Python:
```dockerfile
FROM python:3.11-slim
WORKDIR /app

# Copy requirements
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy application
COPY . .

# Expose port
EXPOSE 8000

# Start command
CMD ["python", "app.py"]
```

Example for Go:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

::: tip Dockerfile Path
If your Dockerfile is not in the root directory, specify the path in the app settings (e.g., `./docker/Dockerfile` or `./backend/Dockerfile`).
:::

### Dockerfile Support

### Dockerfile Support

Mist supports custom Dockerfiles in your repository:

1. **Default location**: `./Dockerfile` in repository root
2. **Custom path**: Specify via `dockerfilePath` setting (e.g., `./docker/Dockerfile.prod`)
3. **Root directory**: Set `rootDirectory` for builds in subdirectories (e.g., `/services/api`)

The build process:
- Uses your specified Dockerfile
- Passes all environment variables as `--build-arg`
- Tags image with commit hash for version tracking
- Stores build logs for debugging

```dockerfile
# Example multi-service monorepo structure
# Repository structure:
# /services
#   /api
#     Dockerfile
#     ...
#   /worker
#     Dockerfile
#     ...

# For API service, set:
# Root Directory: /services/api
# Dockerfile Path: ./Dockerfile

FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --production
COPY . .
EXPOSE 3000
CMD ["node", "server.js"]
```

::: tip Multi-Stage Builds
Use multi-stage Dockerfiles to keep production images small and secure. Copy only necessary artifacts from build stage to runtime stage.
:::

### Health Checks

::: warning Coming Soon
Health checks are an upcoming feature. Container health monitoring is not yet implemented.
:::

Health checks will allow automatic container monitoring:
- HTTP endpoint checks (e.g., `/health`)
- Configurable check intervals and timeouts
- Automatic restart on health check failures
- Status tracking in dashboard

## Git Integration

### Connecting a Repository

For web and service applications:

1. **Install GitHub App** on your GitHub account/organization
2. **Grant access** to repositories in Mist settings
3. **Select repository** when creating application
4. **Choose branch** to deploy from

### Configuration Settings

After creating your application, configure Git integration in the app settings:

- **Git Repository**: Repository name format `owner/repo` (e.g., `myorg/myapp`)
- **Git Branch**: Branch name to deploy (e.g., `main`, `develop`)
- **Root Directory**: Subdirectory for monorepos (optional, default: `/`)
- **Dockerfile Path**: Path to Dockerfile (optional, default: `./Dockerfile`)

::: warning Dockerfile Required
Your repository must contain a valid Dockerfile. Build and start commands are not yet supported - all applications are deployed using Docker.
:::

```javascript
// API example
{
  "appId": 1,
  "gitRepository": "myorg/myapp",
  "gitBranch": "main",
  "rootDirectory": "/",
  "dockerfilePath": "./Dockerfile"
}
```

## Environment Variables

Applications can have unlimited environment variables for configuration. See the [Environment Variables guide](./environment-variables) for details.

### Adding Variables

1. Go to **Environment Variables** tab
2. Click **"Add Variable"** or use **"Bulk Paste"**
3. Enter key-value pairs
4. Variables are applied on next deployment

### Build-time & Runtime

All environment variables are available in both phases:

- **Build-time**: Passed as `--build-arg` during Docker image build
- **Runtime**: Passed as `-e` flags when container starts
- **Both**: All variables are automatically available in both phases

Example Dockerfile using build args:
```dockerfile
FROM node:18-alpine

# Build arg accessible during build
ARG NODE_ENV
ARG API_URL

WORKDIR /app
COPY package*.json ./
RUN npm ci --production=$NODE_ENV

COPY . .

# Runtime environment variables
ENV NODE_ENV=$NODE_ENV
ENV API_URL=$API_URL

CMD ["node", "server.js"]
```

::: warning Secrets Management
Never commit secrets to Git. Always use environment variables for sensitive data like API keys, passwords, and tokens.
:::

## Deployment

### Deployment Process

When you deploy an application, Mist follows this pipeline:

**For Web/Service Apps:**
1. **Repository Clone**: Clones your Git repository to `/var/lib/mist/projects/{projectId}/apps/{appName}`
2. **Image Build**: Builds Docker image using your Dockerfile with environment variables as build args
3. **Container Stop**: Gracefully stops existing container (if running)
4. **Container Start**: Starts new container with updated configuration
5. **Domain Routing**: Updates Traefik routes (for web apps with domains)

**For Database Apps:**
1. **Image Pull**: Pulls the specified Docker image from Docker Hub
2. **Container Stop**: Stops existing container (if running)
3. **Container Start**: Starts new container with template configuration

### Manual Deployment

Trigger deployment from dashboard:

1. Click **"Deploy"** button
2. Latest commit is fetched automatically
3. Deployment added to queue
4. Build starts when worker is available
5. Watch real-time progress and logs

For databases, deployment uses the configured Docker image version from the template.

### Automatic Deployment

With GitHub App configured:

- **Push events**: Auto-deploy on push to configured branch
- **Webhooks**: GitHub triggers deployment automatically
- **Status updates**: Real-time deployment progress

[Learn more about deployments →](./deployments)

## Deployment Strategies

Mist uses a rolling deployment strategy:

- **Rolling** (default): Stops old container, then starts new container
- Zero-downtime deployments coming soon

```javascript
{
  "deploymentStrategy": "rolling"
}
```

## Domains

Add custom domains to web applications:

1. Go to **Domains** tab
2. Click **"Add Domain"**
3. Enter domain name (e.g., `app.example.com`)
4. Configure DNS A record to point to your server
5. Wait for DNS verification
6. SSL certificate issued automatically via Let's Encrypt

[Learn more about domains →](./domains)

## Container Controls

### Networking

All containers are connected to the `traefik-net` Docker bridge network for communication:

**Web Apps with Domains:**
- Connected to `traefik-net`
- Routed through Traefik reverse proxy
- Automatic SSL/TLS via Let's Encrypt
- HTTP to HTTPS redirect enabled
- Accessible via custom domains

**Web Apps without Domains:**
- Exposed directly on host port (e.g., `-p 3000:3000`)
- Accessible via `http://server-ip:port`
- No SSL by default

**Service Apps:**
- Connected to `traefik-net`
- No external port exposure
- Accessible by other containers via container name: `app-{appId}`

**Database Apps:**
- Connected to `traefik-net`
- Internal DNS resolution
- Accessible by other apps via container name: `app-{appId}`
- Example connection: `postgres://user:pass@app-123:5432/dbname`

::: tip Inter-Container Communication
Applications can communicate with each other using the container name format `app-{appId}` where `{appId}` is the application ID shown in the dashboard.
:::

### Start/Stop/Restart

Control container lifecycle from dashboard or API:

- **Stop**: Gracefully stops container
- **Start**: Starts stopped container
- **Restart**: Restarts running container

### Application Status

Container states:

- **`stopped`**: Container is not running
- **`running`**: Container is active and healthy
- **`error`**: Container failed to start or crashed
- **`building`**: Docker image is being built
- **`deploying`**: Container is being deployed

## Monitoring

### Real-time Logs

View live container logs via WebSocket:

- **Real-time streaming**: See logs as they're generated
- **Last 100 lines**: Initial connection shows tail
- **Auto-reconnect**: Maintains connection with ping/pong
- **Timestamps**: Each log line includes timestamp

[View logs tab in dashboard or use WebSocket API](../api/websockets)

### Deployment History

Track all deployments:

- **Deployment number**: Sequential numbering
- **Commit info**: Hash and message (for Git apps)
- **Status**: pending → building → deploying → success/failed
- **Duration**: Build and deployment time
- **Build logs**: Complete build output
- **Active deployment**: Currently running version

### Container Statistics

Monitor basic container information:

- Container name and ID
- Container state (running, stopped, exited)
- Container uptime

::: tip Coming Soon
Advanced metrics including CPU usage, memory consumption, and network I/O are coming soon.
:::

## Application Settings

Update application configuration:

- **General**: Name, description
- **Git**: Repository URL, branch, paths
- **Resources**: CPU, memory limits
- **Container**: Restart policy
- **Status**: Control running state

::: tip Audit Logging
All configuration changes are logged in the audit log with user information and timestamps.
:::

## Rollback

::: warning Coming Soon
Rollback functionality is an upcoming feature. You cannot currently rollback to previous deployments.
:::

When available, you'll be able to:
1. View deployment history
2. Select a previous successful deployment
3. Click "Rollback" to redeploy that version
4. Previous Docker images will be reused (no rebuild needed)

## Service Templates

Database applications use pre-configured service templates:

### How Templates Work

1. **Template Selection**: Choose from available database/service templates
2. **Auto-Configuration**: CPU, memory, port, and environment variables set from template
3. **Docker Image**: Pre-defined image and version from template
4. **One-Click Deploy**: No Dockerfile or Git repository required

### Template Properties

Each template defines:
- `dockerImage`: Docker Hub image name
- `dockerImageVersion`: Image tag (e.g., `latest`, `14-alpine`)
- `defaultPort`: Container port
- `recommendedCpu`: Suggested CPU cores
- `recommendedMemory`: Suggested memory in MB
- `defaultEnvVars`: Pre-configured environment variables
- `volumeRequired`: Whether persistent volume is needed

### Available Templates

Common templates include:
- **PostgreSQL**: Full-featured relational database
- **MySQL**: Popular open-source database
- **Redis**: In-memory data store and cache
- **MongoDB**: Document-oriented database
- **MariaDB**: MySQL-compatible database

::: tip Coming Soon
Custom template creation and management by administrators is an upcoming feature. Currently, only pre-defined templates are available.
:::

## Best Practices

### Application Structure

Organize your repository:

```
my-app/
├── src/               # Source code
├── Dockerfile         # Production Dockerfile
├── .dockerignore      # Exclude unnecessary files
├── package.json       # Dependencies
├── .env.example       # Document required env vars
└── README.md          # Setup instructions
```

### Resource Planning

- **Start small**: Begin with default limits
- **Database memory**: Databases need sufficient RAM
- **CPU limits**: Set realistic limits for consistent performance

::: tip Coming Soon
Once metrics are available, you'll be able to monitor usage and scale accordingly.
:::

### Environment Variables

- Document all required variables in README
- Use descriptive key names (e.g., `DATABASE_URL` not `DB`)
- Separate development and production configs
- Rotate secrets regularly

### Build Optimization

Use `.dockerignore` to exclude files:

```
node_modules
.git
.env
*.log
*.md
.vscode
.idea
```

Use multi-stage builds:

```dockerfile
# Build stage
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production stage
FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
CMD ["node", "dist/server.js"]
```

### Health Check Implementation

Add health check endpoint to your application:

```javascript
// Express.js example
app.get('/health', (req, res) => {
  // Check database connection
  // Check external dependencies
  res.status(200).json({ 
    status: 'healthy',
    timestamp: new Date().toISOString()
  });
});
```

```python
# Flask example
@app.route('/health')
def health():
    return {'status': 'healthy'}, 200
```

## Troubleshooting

### Build Failures

Check deployment logs for errors:

- **Missing dependencies**: Check your package manager files (package.json, requirements.txt, etc.)
- **Dockerfile errors**: Verify Dockerfile syntax and commands
- **Build args missing**: Ensure required environment variables are set
- **Out of disk space**: Clean up old Docker images with `docker system prune`
- **Permission denied**: Check file permissions and COPY commands in Dockerfile
- **Git clone failed**: Verify GitHub App installation and repository access

### Container Won't Start

Common issues:

- **Dockerfile CMD missing**: Container needs a command to run
- **Port mismatch**: Dockerfile EXPOSE port should match app configuration
- **Missing env vars**: Required environment variables not set
- **Application crashes**: Check container logs for runtime errors
- **Resource limits**: Container may be OOM killed if memory limit too low

### Application Not Accessible

Checklist:

- ✅ Container status is `running` (green indicator)
- ✅ Domain has A record pointing to server
- ✅ DNS has propagated (check with `dig` or `nslookup`)
- ✅ Firewall allows ports 80/443
- ✅ Application is listening on correct port
- ✅ No errors in container logs

### Performance Issues

Investigate:

- Check CPU/memory usage in metrics
- Review container logs for errors
- Increase resource limits if needed
- Check database query performance
- Review application profiling data

## API Reference

See the [Applications API documentation](../api/applications) for programmatic access to:

- Create applications
- Update configuration  
- Control containers
- Get status and logs
- Manage resources

## Related Topics

- [Deployments](./deployments) - Deploy your applications
- [Environment Variables](./environment-variables) - Configure applications
- [Domains](./domains) - Add custom domains
- [Git Integration](./git-integration) - Connect GitHub repositories
- [Docker Compose](./docker-compose) - Multi-container deployments
- [Logs](./logs) - Monitor application logs
