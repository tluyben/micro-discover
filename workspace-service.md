# 🚀 Workspace Service API Documentation

This document outlines the API endpoints for the Workspace Service in the micro-discover project.

## 📋 Endpoints

### 1. Create Workspace 🆕

- **URL**: `/workspaces`
- **Method**: `POST`
- **Description**: Creates a new workspace

#### Request Body
```json
{
  "name": "My Workspace",
  "user_id": 1
}
```

#### Response
```json
{
  "id": 1,
  "name": "My Workspace",
  "user_id": 1,
  "subdomain": "abcd1234",
  "ips": ["10.0.0.1"]
}
```

### 2. Get Workspaces 📊

- **URL**: `/workspaces`
- **Method**: `GET`
- **Description**: Retrieves all workspaces

#### Response
```json
[
  {
    "id": 1,
    "name": "My Workspace",
    "user_id": 1,
    "subdomain": "abcd1234",
    "ips": ["10.0.0.1"]
  },
  {
    "id": 2,
    "name": "Another Workspace",
    "user_id": 2,
    "subdomain": "efgh5678",
    "ips": ["10.0.0.2"]
  }
]
```

### 3. Get Workspace 🔍

- **URL**: `/workspaces/{id}`
- **Method**: `GET`
- **Description**: Retrieves a specific workspace by ID

#### Response
```json
{
  "id": 1,
  "name": "My Workspace",
  "user_id": 1,
  "subdomain": "abcd1234",
  "ips": ["10.0.0.1"]
}
```

### 4. Update Workspace 🔄

- **URL**: `/workspaces/{id}`
- **Method**: `PUT`
- **Description**: Updates an existing workspace

#### Request Body
```json
{
  "name": "Updated Workspace Name",
  "user_id": 1
}
```

#### Response
```json
{
  "id": 1,
  "name": "Updated Workspace Name",
  "user_id": 1,
  "subdomain": "abcd1234",
  "ips": ["10.0.0.1"]
}
```

### 5. Delete Workspace 🗑️

- **URL**: `/workspaces/{id}`
- **Method**: `DELETE`
- **Description**: Deletes a workspace

#### Response
- Status: 204 No Content

### 6. Create Workspace Role 👥

- **URL**: `/workspace-roles`
- **Method**: `POST`
- **Description**: Assigns a role to a user for a specific workspace

#### Request Body
```json
{
  "user_id": 1,
  "role": "admin",
  "workspace_id": 1
}
```

#### Response
```json
{
  "id": 1,
  "user_id": 1,
  "role": "admin",
  "workspace_id": 1
}
```

### 7. Get Workspace Roles 👥📊

- **URL**: `/workspace-roles`
- **Method**: `GET`
- **Description**: Retrieves all workspace roles

#### Response
```json
[
  {
    "id": 1,
    "user_id": 1,
    "role": "admin",
    "workspace_id": 1
  },
  {
    "id": 2,
    "user_id": 2,
    "role": "developer",
    "workspace_id": 1
  }
]
```

### 8. Update Workspace Role 🔄👥

- **URL**: `/workspace-roles/{id}`
- **Method**: `PUT`
- **Description**: Updates an existing workspace role

#### Request Body
```json
{
  "user_id": 1,
  "role": "developer",
  "workspace_id": 1
}
```

#### Response
```json
{
  "id": 1,
  "user_id": 1,
  "role": "developer",
  "workspace_id": 1
}
```

### 9. Delete Workspace Role 🗑️👥

- **URL**: `/workspace-roles/{id}`
- **Method**: `DELETE`
- **Description**: Deletes a workspace role

#### Response
- Status: 204 No Content

## 🏗️ Data Models

### Workspace
- `id`: int
- `name`: string
- `user_id`: int
- `subdomain`: string
- `ips`: []string

### WorkspaceRole
- `id`: int
- `user_id`: int
- `role`: string
- `workspace_id`: int

## 🔐 Authentication

All endpoints require authentication. Include the user's JWT token in the Authorization header of each request.

## 🚦 Error Handling

The API uses standard HTTP status codes to indicate the success or failure of requests. In case of an error, the response body will contain a JSON object with an `error` field describing the issue.

Example error response:
```json
{
  "error": "Workspace not found"
}
```

## 📝 Notes

- The `subdomain` field is automatically generated when creating a new workspace.
- The `ips` field is managed by the system and cannot be directly modified by clients.
- Workspace roles determine the permissions a user has within a specific workspace.

This API documentation provides a comprehensive overview of the Workspace Service endpoints, including request/response formats, data models, and important notes for developers integrating with the service.