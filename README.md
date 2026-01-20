# Mist

A lightweight, self-hostable Platform-as-a-Service built for developers. Deploy Docker applications from Git with automatic builds, custom domains, SSL certificates, and real-time monitoring.

[![PR Build Check](https://github.com/corecollectives/Mist/actions/workflows/pr-build-check.yml/badge.svg)](https://github.com/corecollectives/Mist/actions/workflows/pr-build-check.yml)
[![Tests](https://github.com/corecollectives/Mist/actions/workflows/tests.yml/badge.svg)](https://github.com/corecollectives/Mist/actions/workflows/tests.yml)

## Quick Start

Install Mist on your server with a single command:

```bash
curl -sSL https://trymist.cloud/install.sh | bash
```

**Requirements:**
- Linux server (Ubuntu 20.04+ recommended)
- Docker installed
- Root or sudo access

After installation, access the dashboard at `http://your-server-ip:8080`

## Features

-  **Easy Deployment** - Deploy from GitHub with automatic builds
-  **Git Integration** - Webhooks for automatic deployments on push
-  **Domains & SSL** - Custom domains with automatic Let's Encrypt certificates
-  **Real-time Monitoring** - Live logs and system metrics via WebSocket
-  **Database Services** - One-click PostgreSQL, MySQL, Redis, and MongoDB
-  **Secure** - JWT authentication with role-based access control
-  **Ultra Lightweight** - Single Go binary with embedded SQLite (~20MB RAM)

## Documentation

Full documentation is available at [trymist.cloud](https://trymist.cloud/guide/getting-started.html)

## Development

### Setup Development Environment

Install `fyrer`:

```bash
cargo install fyrer
```

Clone and start the development environment:

```bash
git clone https://github.com/corecollectives/mist
cd mist
fyrer
```

## Community

- [GitHub](https://github.com/corecollectives/mist)
- [Discord](https://discord.gg/hr6TCQDDkj)
- [Documentation](https://trymist.cloud/guide/getting-started.html)

## License

MIT License - see [LICENSE](LICENSE) for details
