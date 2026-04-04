# API Overview

Mist provides a comprehensive REST API for programmatic access to all features.

## Base URL

```
https://your-mist-instance.com/api
```

## Authentication

All API requests require authentication using JWT tokens stored in HTTP-only cookies.

### Via Browser

Cookies are set automatically after login.

### Via cURL/API Client

1. Login to get session cookie:

```bash
curl -X POST https://mist.example.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password"}' \
  -c cookies.txt
```

2. Use cookie in subsequent requests:

```bash
curl https://mist.example.com/api/apps \
  -b cookies.txt
```

## API Endpoints

### Authentication
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user

### Projects
- `GET /api/projects` - List projects
- `POST /api/projects/create` - Create project
- `GET /api/projects/:id` - Get project
- `PUT /api/projects/update` - Update project
- `DELETE /api/projects/delete` - Delete project

### Applications
- `POST /api/apps/get` - List applications
- `POST /api/apps/create` - Create application
- `POST /api/apps/update` - Update application
- `POST /api/apps/delete` - Delete application

### Deployments
- `POST /api/deployments/create` - Trigger deployment
- `POST /api/deployments/get` - List deployments
- `GET /api/deployments/:id` - Get deployment details

### Environment Variables
- `POST /api/apps/envs/create` - Create variable
- `POST /api/apps/envs/get` - List variables
- `PUT /api/apps/envs/update` - Update variable
- `DELETE /api/apps/envs/delete` - Delete variable

### Domains
- `POST /api/apps/domains/create` - Add domain
- `POST /api/apps/domains/get` - List domains
- `DELETE /api/apps/domains/delete` - Remove domain

### Users (Admin Only)
- `POST /api/users/get` - List users
- `POST /api/users/create` - Create user
- `DELETE /api/users/delete` - Delete user

### GitHub Integration
- `POST /api/github/repos` - List repositories
- `POST /api/github/branches` - List branches
- `POST /api/webhooks/github` - Webhook endpoint

## WebSocket Endpoints

### Container Logs
```
ws://your-mist-instance.com/api/ws/logs/:containerId
```

### System Metrics
```
ws://your-mist-instance.com/api/ws/metrics
```

### Deployment Status
```
ws://your-mist-instance.com/api/ws/deployment/:deploymentId
```

## Response Format

### Success Response

```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response

```json
{
  "success": false,
  "error": "Error message"
}
```

## Rate Limiting

Currently no rate limiting is enforced. This will be added in future releases.

## API Versioning

The API is currently unversioned. Breaking changes will be avoided when possible.

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 API Improvements</h4>
</div>

- **API Tokens** - Generate tokens for non-browser clients
- **OpenAPI Specification** - Interactive API documentation
- **API Versioning** - v1, v2, etc.
- **Rate Limiting** - Prevent API abuse
- **Webhooks** - Subscribe to events
- **GraphQL API** - Alternative to REST

## Detailed Documentation

- [Authentication API](./authentication)
- [Projects API](./projects)
- [Applications API](./applications)
- [Deployments API](./deployments)
- [Environment Variables API](./environment-variables)
- [Domains API](./domains)
- [Users API](./users)
- [GitHub API](./github)
- [WebSockets API](./websockets)
