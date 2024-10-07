# User Service 👤

The User Service manages user accounts and authentication for the micro-discover platform.

## Endpoints 🛣️

### 1. Create User 🆕

**POST** `/users`

Creates a new user account.

**Request Body:**
```json
{
  "username": "john.doe@example.com",
  "password": "securePassword123",
  "firstName": "John",
  "lastName": "Doe"
}
```

**Response:**
```json
{
  "id": 1,
  "username": "john.doe@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "createdAt": "2023-06-15T10:30:00Z"
}
```

### 2. Get User 🔍

**GET** `/users/{id}`

Retrieves user information by ID.

**Response:**
```json
{
  "id": 1,
  "username": "john.doe@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "createdAt": "2023-06-15T10:30:00Z"
}
```

### 3. Update User ✏️

**PUT** `/users/{id}`

Updates user information.

**Request Body:**
```json
{
  "firstName": "Johnny",
  "lastName": "Doe"
}
```

**Response:**
```json
{
  "id": 1,
  "username": "john.doe@example.com",
  "firstName": "Johnny",
  "lastName": "Doe",
  "updatedAt": "2023-06-16T14:45:00Z"
}
```

### 4. Delete User ❌

**DELETE** `/users/{id}`

Deletes a user account.

**Response:**
```
204 No Content
```

### 5. User Login 🔐

**POST** `/login`

Authenticates a user and returns a JWT token.

**Request Body:**
```json
{
  "username": "john.doe@example.com",
  "password": "securePassword123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresAt": "2023-06-17T10:30:00Z"
}
```

### 6. Get All Users 📋

**GET** `/users`

Retrieves a list of all users (admin only).

**Response:**
```json
[
  {
    "id": 1,
    "username": "john.doe@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "createdAt": "2023-06-15T10:30:00Z"
  },
  {
    "id": 2,
    "username": "jane.smith@example.com",
    "firstName": "Jane",
    "lastName": "Smith",
    "createdAt": "2023-06-16T09:15:00Z"
  }
]
```

## User Model 📊

```go
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Password  string    `json:"-"` // Not returned in JSON responses
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

## Authentication 🔒

All endpoints except `Create User` and `User Login` require a valid JWT token in the `Authorization` header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Error Handling ⚠️

The service returns appropriate HTTP status codes and error messages in JSON format:

```json
{
  "error": "Invalid credentials",
  "code": "AUTH_001"
}
```

Common error codes:
- `AUTH_001`: Authentication failed
- `USER_001`: User not found
- `USER_002`: Invalid input data

## Rate Limiting 🚦

The service implements rate limiting to prevent abuse. Limits are set to 100 requests per minute per IP address.

## Logging and Monitoring 📈

All API calls are logged for auditing purposes. The service integrates with Prometheus for monitoring and alerting.

---

This completes the documentation for the User Service in the micro-discover platform. It includes all endpoints, data models, authentication details, error handling, and other important information for developers working with this service.