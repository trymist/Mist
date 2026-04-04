# WebSockets API

Real-time data streaming via WebSocket connections. Mist provides WebSocket endpoints for live container logs, deployment status, and system metrics.

## Container Logs (Real-time)

`WS /ws/container/logs?appId={appId}`

Streams live container logs from a running application. This provides real-time log output as the container generates it.

### Connection

```javascript
const ws = new WebSocket('wss://mist.example.com/ws/container/logs?appId=1');

ws.onopen = () => {
  console.log('Connected to container logs');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  switch(message.type) {
    case 'status':
      console.log('Container status:', message.data);
      break;
    case 'log':
      console.log('Log:', message.data.line);
      break;
    case 'error':
      console.error('Error:', message.data.message);
      break;
    case 'end':
      console.log('Log stream ended');
      break;
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected from container logs');
};
```

### Message Types

**Status Update (initial connection):**
```json
{
  "type": "status",
  "timestamp": "2025-01-15T10:30:00Z",
  "data": {
    "container": "mist-app-1",
    "state": "running",
    "status": "Up 2 hours"
  }
}
```

**Log Line:**
```json
{
  "type": "log",
  "timestamp": "2025-01-15T10:30:01Z",
  "data": {
    "line": "2025-01-15 10:30:01 [INFO] Server listening on port 3000"
  }
}
```

**Error:**
```json
{
  "type": "error",
  "timestamp": "2025-01-15T10:30:05Z",
  "data": {
    "message": "Container not found"
  }
}
```

**Stream End:**
```json
{
  "type": "end",
  "timestamp": "2025-01-15T10:35:00Z",
  "data": {
    "message": "Log stream ended"
  }
}
```

### Features

- **Auto-reconnect**: Includes ping/pong messages every 30 seconds to keep connection alive
- **Buffer size**: Handles large log lines up to 1MB
- **Tail logs**: Shows last 100 lines on connection
- **Real-time**: Streams logs as they're generated

### Notes

::: warning Container Must Be Running
The WebSocket connection requires the container to be in a `running` state. If the container is stopped, you'll receive an error message.
:::

## Deployment Logs & Status

`WS /ws/deployment/logs?deploymentId={deploymentId}`

Streams real-time deployment progress, build logs, and status updates during a deployment.

### Connection

```javascript
const ws = new WebSocket('wss://mist.example.com/ws/deployment/logs?deploymentId=15');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  if (message.type === 'status') {
    updateProgressBar(message.data.progress);
    updateStage(message.data.stage);
  } else if (message.type === 'log') {
    appendLog(message.data.line);
  } else if (message.type === 'error') {
    showError(message.data.message);
  }
};
```

### Message Types

**Status Update:**
```json
{
  "type": "status",
  "timestamp": "2025-01-15T10:30:15Z",
  "data": {
    "deployment_id": 15,
    "status": "building",
    "stage": "building_image",
    "progress": 50,
    "message": "Building Docker image...",
    "error_message": ""
  }
}
```

**Log Line:**
```json
{
  "type": "log",
  "timestamp": "2025-01-15T10:30:16Z",
  "data": {
    "line": "Step 3/8 : COPY package*.json ./"
  }
}
```

**Error:**
```json
{
  "type": "error",
  "timestamp": "2025-01-15T10:30:30Z",
  "data": {
    "message": "Docker build failed: exit code 1"
  }
}
```

### Status Flow

1. **pending** (0%) - Deployment queued
2. **building** (25-75%) - Building Docker image
3. **deploying** (75-90%) - Starting container
4. **success** (100%) - Deployment complete
5. **failed** (error) - Deployment failed

### Stages

- `cloning_repo` - Cloning Git repository
- `building_image` - Building Docker image
- `stopping_old` - Stopping previous container
- `starting_new` - Starting new container
- `completed` - Deployment finished

## System Metrics

`WS /ws/metrics`

Streams real-time system metrics including CPU usage, memory, disk space, and Docker stats.

### Connection

```javascript
const ws = new WebSocket('wss://mist.example.com/ws/metrics');

ws.onmessage = (event) => {
  const metrics = JSON.parse(event.data);
  updateDashboard(metrics);
};
```

### Message Format

```json
{
  "timestamp": "2025-01-15T10:30:00Z",
  "cpu": {
    "usage": 25.5,
    "cores": 4
  },
  "memory": {
    "used": 4294967296,
    "total": 17179869184,
    "percent": 25.0
  },
  "disk": {
    "used": 53687091200,
    "total": 536870912000,
    "percent": 10.0
  },
  "docker": {
    "containers_running": 5,
    "containers_total": 8,
    "images": 12
  }
}
```

### Update Frequency

Metrics are broadcast every **1 second** to all connected clients.

## React/TypeScript Example

Here's a complete example using React hooks:

```typescript
import { useEffect, useState } from 'react';

interface LogMessage {
  type: 'log' | 'status' | 'error' | 'end';
  timestamp: string;
  data: any;
}

export function useContainerLogs(appId: number) {
  const [logs, setLogs] = useState<string[]>([]);
  const [status, setStatus] = useState<string>('connecting');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const ws = new WebSocket(
      `wss://mist.example.com/ws/container/logs?appId=${appId}`
    );

    ws.onopen = () => {
      setStatus('connected');
    };

    ws.onmessage = (event) => {
      const message: LogMessage = JSON.parse(event.data);

      switch (message.type) {
        case 'log':
          setLogs(prev => [...prev, message.data.line]);
          break;
        case 'status':
          setStatus(message.data.state);
          break;
        case 'error':
          setError(message.data.message);
          break;
        case 'end':
          setStatus('ended');
          break;
      }
    };

    ws.onerror = () => {
      setError('WebSocket connection error');
      setStatus('error');
    };

    ws.onclose = () => {
      setStatus('disconnected');
    };

    return () => {
      ws.close();
    };
  }, [appId]);

  return { logs, status, error };
}
```

## Authentication

WebSocket connections inherit authentication from HTTP cookies. Make sure the user is authenticated before establishing WebSocket connections.

::: tip Secure Connections
In production, always use `wss://` (WebSocket Secure) instead of `ws://` for encrypted connections.
:::

## Error Handling

Common error scenarios:

- **Container not found**: Application doesn't have a running container
- **Permission denied**: User doesn't have access to the application
- **Connection timeout**: Network issues or firewall blocking WebSocket
- **Container not running**: Container is stopped or in error state

## Browser Compatibility

WebSocket API is supported in all modern browsers:
- Chrome/Edge 16+
- Firefox 11+
- Safari 7+
- Opera 12.1+

For the complete implementation, see [server/websockets](https://github.com/trymist/mist/tree/main/server/websockets).
