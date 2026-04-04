# What is Mist?

Mist is a **self-hostable Platform-as-a-Service (PaaS)** that allows you to deploy and manage Docker-based applications on your own infrastructure. Think of it as your own private Heroku, but with more control and zero vendor lock-in.

## Overview

Mist simplifies the deployment workflow by providing:

- **Automated Deployments** - Push to Git and deploy automatically
- **Container Management** - Docker-based isolation for applications
- **Real-time Monitoring** - WebSocket-powered live logs and metrics
- **Team Collaboration** - Projects with multi-user access
- **Domain Management** - Custom domains with automatic SSL and wildcard domain support

## Architecture

Mist is built with a **monolithic Go backend** and a **React + Vite frontend**, designed to run efficiently on a single VPS with minimal resources.

### Core Components

```
┌─────────────────────────────────────────────┐
│           Frontend (React + Vite)            │
│         Modern UI with Real-time Updates     │
└──────────────────┬──────────────────────────┘
                   │ REST API + WebSocket
┌──────────────────▼──────────────────────────┐
│          Backend API (Go)                    │
│  • Authentication & Authorization            │
│  • Project & Application Management          │
│  • Deployment Queue & Workers                │
│  • Real-time Log Streaming                   │
└──────────┬────────────────┬──────────────────┘
           │                │
           │                │
┌──────────▼────────┐  ┌────▼──────────────────┐
│  SQLite Database  │  │   Docker Engine       │
│  • Users          │  │  • Build Images       │
│  • Projects       │  │  • Run Containers     │
│  • Apps           │  │  • Manage Networks    │
│  • Deployments    │  │                       │
└───────────────────┘  └───────────────────────┘
                              │
                    ┌─────────▼─────────┐
                    │   Traefik Proxy    │
                    │  • Route Traffic   │
                    │  • SSL/TLS         │
                    └────────────────────┘
```

[Learn more about architecture →](./architecture)

## Key Features

### Deployment Workflow

1. **Connect Repository** - Link your GitHub repository
2. **Configure Build** - Set build and start commands
3. **Add Environment Variables** - Configure runtime environment
4. **Deploy** - Push to Git or manually trigger deployment
5. **Monitor** - Watch real-time logs and metrics

### Technology Stack

**Backend:**
- Go (1.21+) - Fast, compiled language
- SQLite - Embedded database
- Traefik - Reverse proxy and load balancer
- Docker - Container runtime

**Frontend:**
- React 18 - UI framework
- Vite - Build tool and dev server
- Tailwind CSS - Styling
- Shadcn/ui - Component library

## Why Choose Mist?

### 🪶 Ultra Lightweight

- **Single Binary** - No complex installation process
- **No External Database** - SQLite embedded
- **Minimal Dependencies** - Just Docker required
- **Low Resource Usage** - Runs on basic VPS ($5-10/month)

### ⚡ Real-time Experience

- **WebSocket-First** - Instant feedback for all operations
- **Live Logs** - Stream container and system logs in real-time
- **Live Metrics** - CPU, RAM, and disk usage updated every second
- **Deployment Status** - Watch builds progress in real-time
- **System Monitoring** - View Mist backend logs via journalctl integration

### 🔧 Developer Friendly

- **Simple Setup** - One-line installation script
- **Intuitive UI** - Clean, modern dashboard
- **Comprehensive API** - RESTful endpoints for automation
- **CLI Tool** - (Coming Soon) Deploy from terminal

### 🔐 Secure

- **JWT Authentication** - Secure token-based auth
- **HTTP-Only Cookies** - Protection against XSS
- **Audit Logging** - Track all user actions
- **Role-Based Access** - Admin and user roles
- **Environment Encryption** - (Coming Soon) Secrets at rest

## Use Cases

### Personal Projects

Host your side projects, portfolios, and experimental apps without paying for expensive hosting.

### Small Teams

Collaborate with team members on multiple projects with organized workspaces and access control.

### Staging Environments

Run isolated staging environments alongside production apps for testing before deployment.

### Microservices

Deploy and manage multiple microservices with custom domains and networking.

## Comparison with Other PaaS

| Feature | Mist | Coolify | Dokploy | CapRover | Heroku |
|---------|------|---------|---------|----------|--------|
| **Self-hosted** | ✅ | ✅ | ✅ | ✅ | ❌ |
| **Open Source** | ✅ | ✅ | ✅ | ✅ | ❌ |
| **Real-time Monitoring** | ✅ | ✅ | ✅ | ❌ | ✅ |
| **No External DB** | ✅ (SQLite) | ❌ (Postgres) | ❌ (Postgres) | ❌ (Mongo) | N/A |
| **Go Backend** | ✅ | ❌ (Node) | ❌ (Node) | ❌ (Node) | N/A |
| **Single Binary** | ✅ | ❌ | ❌ | ❌ | N/A |
| **Git Integration** | ✅ GitHub | ✅ Multiple | ✅ Multiple | ✅ Multiple | ✅ Multiple |
| **Managed Databases** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **SSL Automation** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Auto Domain Generation** | ✅ | ❌ | ❌ | ✅ | ✅ |

## Getting Started

Ready to deploy your first application with Mist?

[Get Started →](./getting-started)
