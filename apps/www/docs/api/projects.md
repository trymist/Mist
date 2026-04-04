# Projects API

Manage projects via API.

## List Projects

`GET /api/projects`

Returns all projects accessible to the authenticated user.

## Create Project

`POST /api/projects/create`

**Request:**
```json
{
  "name": "My Project",
  "description": "Project description",
  "tags": "production,web"
}
```

## Get Project

`GET /api/projects/:id`

Returns project details including applications and members.

## Update Project

`PUT /api/projects/update`

**Request:**
```json
{
  "id": 1,
  "name": "Updated Name",
  "description": "New description"
}
```

## Delete Project

`DELETE /api/projects/delete`

**Request:**
```json
{
  "id": 1
}
```

See [server/api/handlers/projects](https://github.com/trymist/mist/tree/main/server/api/handlers/projects) for details.
