# Deployment Rollback <span class="coming-soon-badge">Coming Soon</span>

Instantly revert to previous deployments.

<div class="coming-soon-banner">
  <h4>🚧 Feature in Development</h4>
  <p>Deployment rollback is a high-priority feature for Phase 1.</p>
</div>

## Planned Features

### One-Click Rollback

- Revert to any previous deployment
- Instant container swap
- Zero configuration required
- Preserves deployment history

### Deployment History

- Keep last N deployments (configurable)
- View complete deployment timeline
- Compare deployments side-by-side
- Track rollback events

### Image Management

- Automatic image retention
- Configurable cleanup policy
- Disk space management
- Image size tracking

## How It Will Work

```bash
# In the dashboard:
1. Go to Deployments tab
2. Find previous successful deployment
3. Click "Rollback to this deployment"
4. Confirm rollback
5. Container switches to previous image instantly
```

## Rollback Strategies

### Quick Rollback
- Keep previous container image
- Instant switch (< 10 seconds)
- Previous environment variables restored

### Full Rollback
- Roll back code, config, and environment
- Complete state restoration
- Useful for major issues

## Expected Release

Target: Q1 2025 - Quick Win Priority

[View implementation plan →](https://github.com/trymist/mist/blob/main/roadmap.md#deployment-rollback)
