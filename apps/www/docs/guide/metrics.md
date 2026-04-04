# System Metrics

Monitor your Mist server's resource usage in real-time.

::: tip Application-Level Metrics
Currently, Mist displays system-level metrics only. Per-application metrics (CPU, memory, network usage per container) are coming soon.
:::

## Available Metrics

### CPU Usage

- **Current Usage**: Real-time CPU percentage
- **Per-Core**: Individual core utilization
- **Load Average**: System load over time

### Memory (RAM)

- **Used**: Currently allocated memory
- **Free**: Available memory
- **Total**: Total system memory
- **Percentage**: Memory usage percentage

### Disk Space

- **Used**: Storage space consumed
- **Free**: Available storage
- **Total**: Total disk capacity
- **Percentage**: Disk usage percentage

## Viewing Metrics

Access system metrics from:

1. **Dashboard** - Overview of all metrics
2. **System Status** page - Detailed metrics view

Metrics update every second via WebSocket connection.

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Enhanced Metrics</h4>
</div>

- **Historical Data** - View metrics over time
- **Metric Charts** - Visualize trends with graphs
- **Per-Application Metrics** - Resource usage by app
- **Network Metrics** - Bandwidth monitoring
- **Disk I/O** - Read/write operations
- **Alerts** - Notifications when thresholds exceeded
- **Metric Export** - Prometheus integration
