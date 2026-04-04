# Getting Started

This guide will walk you through installing Mist and deploying your first application.

## Prerequisites

Before installing Mist, ensure you have:

- A Linux server (Ubuntu 20.04+ or Debian 11+ recommended)
- Docker installed and running
- At least 256MB RAM and 2GB disk space
- Root or sudo access
- A domain name (optional, but recommended for production)

## Installation

### Quick Install

The easiest way to install Mist is using the installation script:

```bash
curl -sSL https://trymist.cloud/install.sh | bash
```

This script will:
1. Download the latest Mist binary
2. Set up Traefik reverse proxy
3. Create systemd service for auto-start
4. Configure firewall rules (ports 80, 443, 8080)

[Learn more about installation →](../deployment/installation)

## First-Time Setup

After installation, Mist will be available at `http://your-server-ip:8080`.

### 1. Create Admin Account

On first visit, you'll see the setup page:

1. Enter admin email and password
2. Click "Create Admin Account"
3. You'll be automatically logged in

### 2. Configure Wildcard Domain (Optional)

For automatic domain generation:

1. Go to **Settings** → **System**
2. Enter your wildcard domain (e.g., `example.com`)
3. Configure DNS with a wildcard A record pointing `*.example.com` to your server
4. New web applications will automatically get domains like `{project}-{app}.example.com`

[Learn more about wildcard domains →](./domains#wildcard-domain-configuration)

### 3. Configure GitHub Integration (Optional)

To enable Git deployments:

1. Go to **Settings** → **Git**
2. Follow the instructions to create a GitHub App
3. Install the app on your repositories

[Learn more about GitHub setup →](../deployment/github-app)

## Deploy Your First Application

### Step 1: Create a Project

Projects organize your applications:

1. Click **"New Project"** in the dashboard
2. Enter a name (e.g., "My Portfolio")
3. Add tags (optional)
4. Click **"Create Project"**

### Step 2: Create an Application

1. Open your project
2. Click **"New Application"**
3. Fill in the basic details:
   - **Name**: Your app name
   - **Description**: Brief description of your app
   - **Port**: Application port (e.g., 3000)
4. Click **"Create Application"**
5. Configure additional settings inside the app:
   - **Git Repository**: Select from connected repos
   - **Branch**: Choose branch to deploy
   - **Dockerfile Path**: Path to your Dockerfile (e.g., `./Dockerfile`)

### Step 3: Add Environment Variables

1. Go to the **Environment** tab
2. Click **"Add Variable"** or use **"Bulk Paste"**
3. Add your environment variables:
   ```
   NODE_ENV=production
   DATABASE_URL=your-database-url
   API_KEY=your-api-key
   ```

### Step 4: Deploy

1. Click **"Deploy"** button
2. Watch real-time build logs
3. Wait for deployment to complete
4. Access your app via the generated URL

## Configure Custom Domain

### Step 1: Add Domain

1. Go to **Domains** tab in your application
2. Click **"Add Domain"**
3. Enter your domain (e.g., `app.example.com`)

### Step 2: Configure DNS

Point your domain to your server:

```
Type: A Record
Name: app (or @ for root domain)
Value: YOUR_SERVER_IP
TTL: 3600 (or auto)
```

### Step 3: Wait for DNS Propagation

DNS changes can take 5 minutes to 48 hours to propagate. Check with:

```bash
dig app.example.com
```

::: tip SSL/TLS Certificate
SSL certificates are automatically provisioned using Traefik and Let's Encrypt. Once your DNS is configured and propagated, your application will automatically get an SSL certificate.
:::

## Next Steps

Now that you have Mist running, explore these topics:

- [**Projects**](./projects) - Organize applications
- [**Deployments**](./deployments) - Deployment strategies
- [**Environment Variables**](./environment-variables) - Managing configuration
- [**Monitoring**](./logs) - Container and system logs
- [**Metrics**](./metrics) - System resource monitoring
- [**Git Integration**](./git-integration) - Auto-deploy on push

## Getting Help

- [GitHub Issues](https://github.com/trymist/mist/issues) - Report bugs
- [GitHub Discussions](https://github.com/trymist/mist/discussions) - Ask questions
- [Documentation](../) - Read more guides

## What's Next?

### Coming Soon Features

- 🚧 Email notifications for deployments
- 🚧 CLI tool for terminal deployments
