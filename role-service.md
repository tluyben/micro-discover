# 🚀 Role Service

The Role Service manages user roles and permissions within the micro-discover system.

## 🔑 Endpoints

### 1. Create Role

**POST** `/roles`

Creates a new role.

**Request Body:**
```json
{
  "name": "developer",
  "permissions": ["read", "write", "execute"]
}
```

**Response:**
```json
{
  "id": 1,
  "name": "developer",
  "permissions": ["read", "write", "execute"],
  "created_at": "2023-06-15T10:30:00Z"
}
```

### 2. Get Role

**GET** `/roles/{id}`

Retrieves a specific role by ID.

**Response:**
```json
{
  "id": 1,
  "name": "developer",
  "permissions": ["read", "write", "execute"],
  "created_at": "2023-06-15T10:30:00Z",
  "updated_at": "2023-06-15T10:30:00Z"
}
```

### 3. Update Role

**PUT** `/roles/{id}`

Updates an existing role.

**Request Body:**
```json
{
  "name": "senior_developer",
  "permissions": ["read", "write", "execute", "deploy"]
}
```

**Response:**
```json
{
  "id": 1,
  "name": "senior_developer",
  "permissions": ["read", "write", "execute", "deploy"],
  "updated_at": "2023-06-15T11:00:00Z"
}
```

### 4. Delete Role

**DELETE** `/roles/{id}`

Deletes a role.

**Response:** HTTP 204 No Content

### 5. List Roles

**GET** `/roles`

Retrieves a list of all roles.

**Response:**
```json
[
  {
    "id": 1,
    "name": "senior_developer",
    "permissions": ["read", "write", "execute", "deploy"]
  },
  {
    "id": 2,
    "name": "junior_developer",
    "permissions": ["read", "write"]
  }
]
```

### 6. Assign Role to User

**POST** `/users/{userId}/roles`

Assigns a role to a user.

**Request Body:**
```json
{
  "role_id": 1
}
```

**Response:**
```json
{
  "user_id": 123,
  "role_id": 1,
  "assigned_at": "2023-06-15T12:00:00Z"
}
```

### 7. Remove Role from User

**DELETE** `/users/{userId}/roles/{roleId}`

Removes a role from a user.

**Response:** HTTP 204 No Content

## 📊 Data Models

### Role
- `id`: int
- `name`: string
- `permissions`: string[]
- `created_at`: datetime
- `updated_at`: datetime

### UserRole
- `user_id`: int
- `role_id`: int
- `assigned_at`: datetime

## 🔒 Authentication

All endpoints require a valid JWT token in the Authorization header.

## 🚦 Rate Limiting

API requests are limited to 100 requests per minute per API key.

## 📝 Notes

- Role names must be unique.
- Deleting a role will automatically remove it from all users who have been assigned that role.
- The 'admin' role is reserved and cannot be modified or deleted through the API.

🎉 Happy role managing!