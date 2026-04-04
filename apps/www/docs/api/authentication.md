# Authentication API

Manage user authentication via API.

## Login

`POST /api/auth/login`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "role": "admin"
    }
  }
}
```

Sets HTTP-only cookie with JWT token.

## Logout

`POST /api/auth/logout`

Clears authentication cookie.

## Get Current User

`GET /api/auth/me`

Returns currently authenticated user.

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "role": "admin",
    "createdAt": "2025-01-15T10:00:00Z"
  }
}
```

See [server/api/handlers/auth](https://github.com/trymist/mist/tree/main/server/api/handlers/auth) for implementation details.
