# Environment Variables API

Manage environment variables via API.

## Create Variable

`POST /api/apps/envs/create`

### Request Body

```json
{
  "appId": 1,
  "key": "API_KEY",
  "value": "secret-value"
}
```

### Response

```json
{
  "success": true,
  "data": {
    "id": 1,
    "appId": 1,
    "key": "API_KEY",
    "value": "secret-value",
    "createdAt": "2025-01-15T10:30:00Z"
  }
}
```

## List Variables

`POST /api/apps/envs/get`

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
  "data": [
    {
      "id": 1,
      "key": "API_KEY",
      "value": "secret-value",
      "createdAt": "2025-01-15T10:30:00Z"
    }
  ]
}
```

## Update Variable

`PUT /api/apps/envs/update`

### Request Body

```json
{
  "id": 1,
  "key": "API_KEY",
  "value": "new-secret-value"
}
```

## Delete Variable

`DELETE /api/apps/envs/delete`

### Request Body

```json
{
  "id": 1
}
```

## Example: Bulk Import

```bash
#!/bin/bash

APP_ID=1

# Read .env file and create variables
while IFS='=' read -r key value; do
  # Skip comments and empty lines
  [[ $key =~ ^#.*$ ]] && continue
  [[ -z $key ]] && continue
  
  # Create variable
  curl -X POST https://mist.example.com/api/apps/envs/create \
    -H "Content-Type: application/json" \
    -b cookies.txt \
    -d "{\"appId\": $APP_ID, \"key\": \"$key\", \"value\": \"$value\"}"
done < .env
```

For complete API reference, see [server/api/handlers/applications/envs.go](https://github.com/trymist/mist/blob/main/server/api/handlers/applications/envs.go).
