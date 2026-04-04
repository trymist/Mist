# Docker Compose Applications

Deploy complex multi-container applications using Docker Compose files.

## Overview

Docker Compose apps allow you to deploy applications defined in a `docker-compose.yml` file. This is ideal for:

- Multi-service applications (web app + database + cache)
- Applications with complex networking requirements
- Services that need to share volumes
- Development environments with multiple containers

## Creating a Compose Application

### Via Dashboard

1. Navigate to your project
2. Click **"New Application"**
3. Select **"Docker Compose"** as the application type
4. Enter application name and description
5. Click **Create Application**

### Configure Git Source

After creating the app:

1. Go to the **Sources** tab
2. Choose **GitHub** (with GitHub App) or **Public Git** (any public repo)
3. Enter repository URL containing your `docker-compose.yml`
4. Select the branch to deploy
5. Save configuration

## docker-compose.yml Requirements

Your repository must contain a `docker-compose.yml` file in the root directory (or the configured root directory).

### Basic Example

```yaml
version: '3.8'

services:
  web:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=postgres://db:5432/myapp
    depends_on:
      - db

  db:
    image: postgres:15
    environment:
      - POSTGRES_PASSWORD=secret
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### Supported Features

- **Build from Dockerfile** - Use `build: .` to build images
- **Pre-built images** - Use `image:` to pull from registries
- **Environment variables** - Defined in compose file or via Mist env vars
- **Volumes** - Named volumes and bind mounts
- **Networking** - Default network created for all services
- **Depends on** - Service startup ordering

::: tip Environment Variables
Environment variables added via Mist dashboard are passed to `docker compose up` and available to all services.
:::

## Service Status

Compose apps show aggregate status for all services:

| Status | Description |
|--------|-------------|
| **Running** | All services are running |
| **Partial** | Some services running, some stopped |
| **Stopped** | All services stopped |

### Individual Container Statistics

The dashboard displays detailed information for each service in your compose stack:

| Field | Description |
|-------|-------------|
| **Service Name** | Name defined in docker-compose.yml |
| **State** | Current state (running, stopped, restarting, exited) |
| **Status** | Container status with uptime (e.g., "Up 2 hours") |

Each container shows its individual state, making it easy to identify which services are running and which may have issues.

## Container Operations

### Start All Services

Runs `docker compose up -d` to start all services in detached mode.

### Stop All Services

Runs `docker compose down` to stop and remove all containers.

### Restart All Services

Runs `docker compose restart` to restart all running containers.

## Logs

View aggregated logs from all services in your compose stack. Logs are combined from all containers with service name prefixes.

## Deployments

When you trigger a deployment:

1. Repository is cloned/updated
2. `docker compose up -d` runs with environment variables
3. Services are started in dependency order
4. Status updates to reflect running services

## Best Practices

### Use Named Volumes

Persist data across deployments:

```yaml
volumes:
  db_data:

services:
  db:
    volumes:
      - db_data:/var/lib/postgresql/data
```

### Set Resource Limits

Prevent runaway containers:

```yaml
services:
  web:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
```

### Health Checks

Add health checks for reliable service discovery:

```yaml
services:
  web:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Limitations

- **No domain routing** - Compose apps don't automatically get Traefik routing. Configure Traefik labels in your compose file if needed.
- **Single stack** - Each Mist app manages one compose stack
- **Root directory** - The `docker-compose.yml` must be in the configured root directory

## Related

- [Git Integration](./git-integration) - Configure repository source
- [Environment Variables](./environment-variables) - Pass configuration to services
- [Applications Guide](./applications) - Overview of all application types
