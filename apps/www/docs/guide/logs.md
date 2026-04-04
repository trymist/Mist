# Logs & Monitoring

Monitor your applications and system with real-time logs and system metrics.

## System Logs

View live logs from the Mist backend server:

1. Navigate to **Logs** page in the main navigation
2. System logs stream in real-time
3. Powered by `journalctl` under the hood

### What's Included

System logs show all output from the Mist Go backend, including:
- Application lifecycle events (startup, shutdown)
- Deployment processing
- API requests and responses
- Container operations
- Error messages and stack traces
- System events and notifications

### Features

- **Real-time streaming** - Logs appear instantly via WebSocket
- **Live updates** - See backend activity as it happens
- **journalctl integration** - Leverages systemd journal for reliable log collection
- **Filterable** - Search and filter system logs
- **Persistent** - Logs stored in systemd journal

::: tip
System logs are useful for debugging deployment issues, monitoring backend activity, and troubleshooting system-level problems.
:::

## Container Logs

View live logs from your running containers:

1. Navigate to your application
2. Go to **Logs** tab
3. Logs stream in real-time via WebSocket

### Features

- **Real-time streaming** - Logs appear instantly
- **Auto-scroll** - Automatically scrolls to latest logs
- **Color coding** - Error, warning, and info levels highlighted
- **Search** - Filter logs by keyword (coming soon)
- **Download** - Export logs to file (coming soon)

## Build Logs

View logs from deployment builds:

1. Go to **Deployments** tab
2. Click on a deployment
3. View complete build output

Build logs show:
- Git clone process
- Dependency installation
- Build command execution
- Docker image creation
- Container startup

## System Metrics

Monitor server resource usage:

### CPU Usage

Real-time CPU utilization displayed as percentage.

### Memory Usage

Current RAM usage and available memory.

### Disk Usage

Storage space used and available on server.

### Update Frequency

Metrics update every second via WebSocket.

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Advanced Monitoring Features</h4>
</div>

- **Log Aggregation** - Search across all logs
- **Log Retention** - Configure how long to keep logs
- **Log Export** - Download logs as files
- **Application Metrics** - Request count, response times
- **Alerting** - Email/Slack notifications for errors
- **Custom Dashboards** - Create metric visualizations
- **Error Tracking** - Sentry-like error monitoring
- **Uptime Monitoring** - Track application availability

## Log Best Practices

### Structured Logging

Use JSON format for easier parsing:

```javascript
console.log(JSON.stringify({
  level: 'info',
  message: 'User logged in',
  userId: 123,
  timestamp: new Date().toISOString()
}));
```

### Log Levels

Use appropriate log levels:

```javascript
console.log('Info message');    // Info
console.warn('Warning message'); // Warning
console.error('Error message');  // Error
```

### Avoid Sensitive Data

Never log:
- Passwords
- API keys
- Credit card numbers
- Personal identification
