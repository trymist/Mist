# Applications API

Manage applications programmatically. All endpoints require authentication via JWT cookie.

## Create Application

`POST /api/apps/create`

Creates a new application in a project. Supports four app types: `web`, `service`, `database`, and `compose`.

### Request Body

```json
{
  "name": "my-app",
  "description": "My application description",
  "projectId": 1,
  "appType": "web",
  "port": 3000,
  "envVars": {
    "NODE_ENV": "production",
    "API_KEY": "secret"
  }
}
```

**For database apps**, use:

```json
{
  "name": "my-postgres",
  "description": "PostgreSQL database",
  "projectId": 1,
  "appType": "database",
  "templateName": "postgresql",
  "envVars": {
    "POSTGRES_PASSWORD": "secret123"
  }
}
```

**For compose apps**, use:

```json
{
  "name": "my-stack",
  "description": "Multi-container application",
  "projectId": 1,
  "appType": "compose"
}
```

### Parameters

- `name` (required): Application name
- `description` (optional): Application description
- `projectId` (required): ID of the project
- `appType` (optional): One of `web`, `service`, `database`, or `compose` (defaults to `web`)
- `port` (optional): Port number for web apps (defaults to 3000)
- `templateName` (required for database apps): Template name (e.g., `postgresql`, `mysql`, `redis`, `mongodb`)
- `envVars` (optional): Key-value map of environment variables

### Response

```json
{
  "success": true,
  "message": "Application created successfully",
  "data": {
    "id": 1,
    "name": "my-app",
    "description": "My application description",
    "project_id": 1,
    "app_type": "web",
    "port": 3000,
    "status": "stopped",
    "created_by": 1,
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

::: tip
For web apps, Mist automatically generates a domain if wildcard DNS is configured.
:::

## List Applications

`POST /api/apps/get`

Retrieves all applications in a project that the authenticated user has access to.

### Request Body

```json
{
  "projectId": 1
}
```

### Response

```json
{
  "success": true,
  "message": "Applications retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "my-app",
      "app_type": "web",
      "status": "running",
      "port": 3000,
      "git_repository": "https://github.com/user/repo",
      "git_branch": "main",
      "created_at": "2025-01-15T10:30:00Z"
    }
  ]
}
```

## Get Application by ID

`POST /api/apps/get/id`

Retrieves a single application by its ID.

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
  "message": "Application retrieved successfully",
  "data": {
    "id": 1,
    "name": "my-app",
    "description": "My application",
    "project_id": 1,
    "app_type": "web",
    "port": 3000,
    "status": "running",
    "git_repository": "https://github.com/user/repo",
    "git_branch": "main",
    "root_directory": "",
    "dockerfile_path": null,
    "deployment_strategy": "rolling",
    "restart_policy": "unless-stopped",
    "cpu_limit": 1.0,
    "memory_limit": 512,
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T12:00:00Z"
  }
}
```

## Update Application

`POST /api/apps/update`

Updates an application's configuration. Only the application owner can update settings.

### Request Body

All fields are optional. Only include fields you want to update:

```json
{
  "appId": 1,
  "name": "updated-app",
  "description": "Updated description",
  "gitRepository": "https://github.com/user/new-repo",
  "gitBranch": "develop",
  "port": 8080,
  "rootDirectory": "/app",
  "dockerfilePath": "docker/Dockerfile",
  "deploymentStrategy": "rolling",
  "status": "running",
  "cpuLimit": 2.0,
  "memoryLimit": 1024,
  "restartPolicy": "always"
}
```

### Parameters

- `appId` (required): Application ID to update
- `name`: Application name
- `description`: Application description
- `gitProviderId`: ID of linked git provider (null for public git)
- `gitRepository`: Git repository name (e.g., `user/repo` for GitHub)
- `gitBranch`: Git branch name
- `gitCloneUrl`: Full clone URL for public git repositories (e.g., `https://github.com/user/repo.git`)
- `port`: Application port
- `rootDirectory`: Root directory for builds
- `dockerfilePath`: Path to Dockerfile
- `deploymentStrategy`: Deployment strategy
- `status`: Application status (`stopped`, `running`, `error`, `building`, `deploying`)
- `cpuLimit`: CPU limit in cores (e.g., 1.0, 2.0)
- `memoryLimit`: Memory limit in MB (e.g., 512, 1024)
- `restartPolicy`: One of `no`, `always`, `on-failure`, `unless-stopped`

::: tip Public Git URL
Use `gitCloneUrl` to deploy from any public repository without setting up a GitHub App. Set `gitProviderId` to `null` when using public git.
:::

### Response

```json
{
  "success": true,
  "message": "Application updated successfully",
  "data": {
    "id": 1,
    "name": "updated-app",
    "git_branch": "develop",
    "updated_at": "2025-01-15T14:00:00Z"
  }
}
```

## Container Controls

### Stop Container

`POST /api/apps/container/stop?appId=1`

Stops a running container.

### Start Container

`POST /api/apps/container/start?appId=1`

Starts a stopped container.

### Restart Container

`POST /api/apps/container/restart?appId=1`

Restarts a container.

### Response (all control endpoints)

```json
{
  "success": true,
  "message": "Container stopped successfully",
  "data": {
    "message": "Container stopped successfully"
  }
}
```

## Get Container Status

`GET /api/apps/container/status?appId=1`

Retrieves the current status of a container.

### Response

```json
{
  "success": true,
  "message": "Container status retrieved successfully",
  "data": {
    "state": "running",
    "status": "Up 2 hours",
    "id": "abc123def456"
  }
}
```

## Get Container Logs

`GET /api/apps/container/logs?appId=1&tail=100`

Retrieves container logs.

### Query Parameters

- `appId` (required): Application ID
- `tail` (optional): Number of lines to retrieve (default: 100)

### Response

```json
{
  "success": true,
  "message": "Container logs retrieved successfully",
  "data": {
    "logs": "2025-01-15T10:30:00Z Server started\n2025-01-15T10:30:01Z Listening on port 3000"
  }
}
```

## cURL Examples

### Create a Web Application

```bash
curl -X POST https://mist.example.com/api/apps/create \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "name": "my-web-app",
    "description": "My Node.js application",
    "projectId": 1,
    "appType": "web",
    "port": 3000,
    "envVars": {
      "NODE_ENV": "production"
    }
  }'
```

### Create a Database Application

```bash
curl -X POST https://mist.example.com/api/apps/create \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "name": "my-postgres",
    "description": "PostgreSQL database",
    "projectId": 1,
    "appType": "database",
    "templateName": "postgresql",
    "envVars": {
      "POSTGRES_PASSWORD": "mysecretpassword",
      "POSTGRES_DB": "mydb"
    }
  }'
```

### Update Application

```bash
curl -X POST https://mist.example.com/api/apps/update \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "appId": 1,
    "name": "updated-app",
    "gitBranch": "develop",
    "cpuLimit": 2.0,
    "memoryLimit": 1024
  }'
```

### Stop Container

```bash
curl -X POST 'https://mist.example.com/api/apps/container/stop?appId=1' \
  -b cookies.txt
```

### Get Container Logs

```bash
curl -X GET 'https://mist.example.com/api/apps/container/logs?appId=1&tail=50' \
  -b cookies.txt
```

::: tip Authentication
All endpoints require authentication via JWT cookie. Use `-b cookies.txt` with cURL to persist cookies after login.
:::

## Related Endpoints

- [Environment Variables API](./environment-variables) - Manage app environment variables
- [Domains API](./domains) - Manage app domains
- [Deployments API](./deployments) - Deploy applications
- [WebSocket API](./websockets) - Real-time container logs

For the complete implementation, see [server/api/handlers/applications](https://github.com/trymist/mist/tree/main/server/api/handlers/applications).
