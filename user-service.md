# 👤 User Service API Documentation

The User Service provides endpoints for managing user accounts in the micro-discover platform.

## 📋 Endpoints

### 🆕 Create User

- **URL**: `/users`
- **Method**: `POST`
- **Description**: Create a new user account

#### Request Body

```json
{
  "username": "user@example.com",
  "password": "securepassword123"
}
```

#### Response

```json
{
  "id": 1,
  "username": "user@example.com"
}
```

### 📖 Get Users

- **URL**: `/users`
- **Method**: `GET`
- **Description**: Retrieve a list of all users

#### Response

```json
[
  {
    "id": 1,
    "username": "user1@example.com"
  },
  {
    "id": 2,
    "username": "user2@example.com"
  }
]
```

### 🔍 Get User

- **URL**: `/users/{id}`
- **Method**: `GET`
- **Description**: Retrieve a specific user by ID

#### Response

```json
{
  "id": 1,
  "username": "user@example.com"
}
```

### 🔄 Update User

- **URL**: `/users/{id}`
- **Method**: `PUT`
- **Description**: Update an existing user's information

#### Request Body

```json
{
  "username": "updateduser@example.com",
  "password": "newpassword123"
}
```

#### Response

```json
{
  "id": 1,
  "username": "updateduser@example.com"
}
```

### ❌ Delete User

- **URL**: `/users/{id}`
- **Method**: `DELETE`
- **Description**: Delete a user account

#### Response

- Status: 204 No Content

## 📝 Notes

- All endpoints return appropriate HTTP status codes (200 for success, 400 for bad requests, 404 for not found, etc.)
- Passwords are hashed before storing in the database
- The service uses email addresses as usernames
- User IDs are automatically generated upon creation
- Password is never returned in responses for security reasons

## 🔐 Authentication

Currently, the User Service does not implement authentication. Future versions may include token-based authentication for secure access to user data.

## 🛠 Error Handling

Errors are returned with appropriate HTTP status codes and error messages in the response body. For example:

```json
{
  "error": "Invalid email address"
}
```

Common error scenarios include:
- Invalid email format
- Duplicate email addresses
- User not found
- Internal server errors

Always check the HTTP status code and response body for error details when interacting with the API.
