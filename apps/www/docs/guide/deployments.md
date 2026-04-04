# Deployments

Learn how to deploy and manage your applications.

## Deployment Process

When you trigger a deployment, Mist:

1. **Clones** your Git repository (for web/service apps)
2. **Builds** Docker image using your Dockerfile with environment variables as build args
3. **Stops** the previous container (if exists)
4. **Starts** new container with updated code
5. **Updates** Traefik routing configuration (for web apps with domains)
6. **Streams** real-time logs to your dashboard

For database applications, Mist pulls the Docker image instead of building from source.

## Triggering Deployments

### Manual Deployment

Click the **"Deploy"** button in your application dashboard.

### Automatic Deployment

Configure GitHub webhooks to auto-deploy on push:

1. Install Mist GitHub App on your repository
2. Enable webhooks in application settings
3. Push to your configured branch
4. Deployment triggers automatically

[Learn more about Git integration →](./git-integration)

## Deployment Status

Monitor deployment progress in real-time:

- 🟡 **Queued** - Waiting in deployment queue
- 🔵 **Building** - Docker image being built
- 🟢 **Running** - Container started successfully
- 🔴 **Failed** - Deployment encountered an error

## Build Logs

View detailed build logs during and after deployment:

1. Go to **Deployments** tab
2. Click on a deployment
3. View streaming logs in real-time

Logs include:
- Git clone output
- Build command execution
- Docker image creation
- Container startup

## Deployment History

Track all deployments with full history:

- Deployment timestamp
- Commit hash and message
- Build duration
- Status (success/failed)
- Triggering user or webhook

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Deployment Features</h4>
</div>

- **Rollback** - Revert to previous deployments with one click
- **Blue-Green Deployments** - Zero-downtime deployments
- **Canary Releases** - Gradual traffic shifting
- **Deployment Hooks** - Pre/post deploy scripts
- **Manual Approval** - Require approval before production deploys
- **Scheduled Deployments** - Deploy at specific times
- **Preview Environments** - Auto-deploy pull requests

## Troubleshooting

### Build Failures

Common causes:
- Missing dependencies
- Incorrect build command
- Dockerfile errors
- Insufficient permissions

**Solution**: Check build logs for specific error messages

### Container Won't Start

Possible issues:
- Wrong start command
- Port configuration mismatch
- Missing environment variables
- Application crashes on startup

**Solution**: Review container logs for errors

### Slow Builds

Optimization tips:
- Use `.dockerignore` to exclude unnecessary files
- Implement Docker layer caching
- Use multi-stage builds
- Optimize dependency installation

## Best Practices

- ✅ Test locally before deploying
- ✅ Use feature branches for experimentation
- ✅ Monitor logs during deployment
- ✅ Set up health checks
- ✅ Keep deployment logs for debugging
