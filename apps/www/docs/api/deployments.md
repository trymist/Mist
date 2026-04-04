# Deployments API

Trigger and manage deployments programmatically. All endpoints require authentication.

## Trigger Deployment

`POST /api/deployment/add`

Triggers a new deployment for an application. For web/service apps, it fetches the latest commit from the configured Git repository. For database apps, it deploys using the service template's Docker image.

### Request Body

```json
{
  "appId": 1
}
```

### Parameters

- `appId` (required): The ID of the application to deploy

### Response

```json
{
  "id": 15,
  "app_id": 1,
  "commit_hash": "a1b2c3d4e5f6",
  "commit_message": "Fix user authentication bug",
  "status": "pending",
  "stage": "",
  "progress": 0,
  "created_at": "2025-01-15T10:30:00Z"
}
```

::: tip Deployment Queue
Deployments are added to a queue and processed by background workers. The deployment status updates from `pending` → `building` → `deploying` → `success` (or `failed`).
:::

## List Deployments

`POST /api/deployment/getByAppId`

Retrieves deployment history for an application, including all previous deployments with their statuses.

### Request Body

```json
{
  "appId": 1
}
```

### Response

```json
{
  "success": true,
  "message": "Deployments retrieved successfully",
  "data": [
    {
      "id": 15,
      "app_id": 1,
      "deployment_number": 15,
      "commit_hash": "a1b2c3d4e5f6",
      "commit_message": "Fix user authentication bug",
      "status": "success",
      "stage": "completed",
      "progress": 100,
      "duration": 45,
      "is_active": true,
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:31:00Z"
    },
    {
      "id": 14,
      "app_id": 1,
      "deployment_number": 14,
      "commit_hash": "f6e5d4c3b2a1",
      "commit_message": "Add user profile page",
      "status": "success",
      "stage": "completed",
      "progress": 100,
      "duration": 38,
      "is_active": false,
      "rolled_back_from": 15,
      "created_at": "2025-01-14T16:20:00Z",
      "updated_at": "2025-01-14T16:21:00Z"
    }
  ]
}
```

### Deployment Statuses

- `pending`: Deployment is queued and waiting to start
- `building`: Application is being built (cloning repo, building Docker image)
- `deploying`: Container is being deployed
- `success`: Deployment completed successfully
- `failed`: Deployment failed (check error message and logs)
- `stopped`: Deployment was manually stopped
- `rolled_back`: Deployment was rolled back to a previous version

## Get Deployment Logs

`GET /api/deployment/logs/:deploymentId`

Retrieves build logs for a specific deployment. These are the logs generated during the build process.

### Response

```
Cloning repository...
Cloning into '/tmp/mist-build-123'...
Building Docker image...
Step 1/8 : FROM node:18-alpine
 ---> 1a2b3c4d5e6f
Step 2/8 : WORKDIR /app
 ---> Running in 7f8e9d0c1b2a
 ---> 3c4d5e6f7a8b
...
Successfully built abc123def456
Successfully tagged mist-app-1:latest
Deployment completed successfully
```

## WebSocket: Watch Deployment Status

`WS /ws/deployment/logs?deploymentId=15`

Connect to a WebSocket to receive real-time deployment status updates and logs.

### Connection

```javascript
const ws = new WebSocket('wss://mist.example.com/ws/deployment/logs?deploymentId=15');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  if (data.type === 'status') {
    console.log('Status:', data.data.status);
    console.log('Progress:', data.data.progress + '%');
    console.log('Stage:', data.data.stage);
  } else if (data.type === 'log') {
    console.log('Log:', data.data.line);
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
    "stage": "cloning_repo",
    "progress": 25,
    "message": "Cloning repository...",
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
    "line": "Step 1/8 : FROM node:18-alpine"
  }
}
```

**Error:**
```json
{
  "type": "error",
  "timestamp": "2025-01-15T10:30:30Z",
  "data": {
    "message": "Failed to build Docker image: exit code 1"
  }
}
```

## Deployment Stages

During a deployment, the system goes through these stages:

1. **cloning_repo** - Cloning Git repository
2. **building_image** - Building Docker image
3. **stopping_old** - Stopping previous container
4. **starting_new** - Starting new container
5. **completed** - Deployment finished successfully

## cURL Examples

### Trigger Deployment

```bash
curl -X POST https://mist.example.com/api/deployment/add \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"appId": 1}'
```

### List Deployments

```bash
curl -X POST https://mist.example.com/api/deployment/getByAppId \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"appId": 1}'
```

### Get Deployment Logs

```bash
curl https://mist.example.com/api/deployment/logs/15 \
  -b cookies.txt
```

::: tip Real-time Updates
For real-time deployment progress, use the WebSocket endpoint. The REST API logs endpoint only returns completed logs after the deployment finishes.
:::

## Related Endpoints

- [Applications API](./applications) - Manage applications
- [WebSocket API](./websockets) - Real-time updates
- [GitHub API](./github) - GitHub integration

For the complete implementation, see [server/api/handlers/deployments](https://github.com/trymist/mist/tree/main/server/api/handlers/deployments) and [server/docker/deployer.go](https://github.com/trymist/mist/blob/main/server/docker/deployer.go).
