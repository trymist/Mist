# GitHub API

GitHub integration endpoints.

## List Repositories

`POST /api/github/repos`

Returns repositories accessible via installed GitHub App.

**Request:**
```json
{
  "installationId": 12345
}
```

## List Branches

`POST /api/github/branches`

**Request:**
```json
{
  "repoFullName": "user/repo",
  "installationId": 12345
}
```

Returns branches for the specified repository.

## Webhook Endpoint

`POST /api/webhooks/github`

Receives GitHub webhook events for auto-deployment.

See [server/api/handlers/github](https://github.com/trymist/mist/tree/main/server/api/handlers/github) for implementation.
