# Users API

Manage users (Admin only).

## List Users

`POST /api/users/get`

Returns all users in the system.

## Create User

`POST /api/users/create`

**Request:**
```json
{
  "email": "newuser@example.com",
  "password": "securepassword",
  "role": "user"
}
```

## Delete User

`DELETE /api/users/delete`

**Request:**
```json
{
  "id": 1
}
```

::: warning Admin Only
These endpoints require admin role.
:::

See [server/api/handlers/auth/users.go](https://github.com/trymist/mist/blob/main/server/api/handlers/auth/users.go) for details.
