# Domains API

Manage application domains.

## Add Domain

`POST /api/apps/domains/create`

**Request:**
```json
{
  "appId": 1,
  "domain": "app.example.com"
}
```

## List Domains

`POST /api/apps/domains/get`

**Request:**
```json
{
  "appId": 1
}
```

## Delete Domain

`DELETE /api/apps/domains/delete`

**Request:**
```json
{
  "id": 1
}
```

See [server/api/handlers/applications/domains.go](https://github.com/trymist/mist/blob/main/server/api/handlers/applications/domains.go) for implementation.
