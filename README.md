# micro-discover

micro-discover is a Go-based microservices discovery and management system. It provides a robust API for managing users, workspaces, applications, and role-based access control.

## Features

- User Management
- Workspace Management
- Application Management
- Role-Based Access Control
- IP Pool Management

## API Documentation

### Users

#### Create User

```
POST /users
Body: {
  "username": "user@example.com",
  "password": "password123"
}
```

#### Get Users

```
GET /users
```

#### Get User

```
GET /users/{id}
```

#### Update User

```
PUT /users/{id}
Body: {
  "username": "updateduser@example.com",
  "password": "newpassword123"
}
```

#### Delete User

```
DELETE /users/{id}
```

### Workspaces

#### Create Workspace

```
POST /workspaces
Body: {
  "name": "My Workspace",
  "user_id": 1
}
```

#### Get Workspaces

```
GET /workspaces
```

#### Get Workspace

```
GET /workspaces/{id}
```

#### Update Workspace

```
PUT /workspaces/{id}
Body: {
  "name": "Updated Workspace",
  "user_id": 2
}
```

#### Delete Workspace

```
DELETE /workspaces/{id}
```

### Applications

#### Create Application

```
POST /apps
Body: {
  "name": "My App",
  "description": "Description of my app",
  "git_hash": "abc123",
  "ip_port": "10.0.0.1:8080",
  "endpoint": "/api",
  "version": "1.0",
  "workspace_id": 1
}
```

#### Get Applications

```
GET /apps
```

#### Get Application

```
GET /apps/{id}
```

#### Update Application

```
PUT /apps/{id}
Body: {
  "name": "Updated App",
  "description": "Updated description",
  "git_hash": "def456",
  "ip_port": "10.0.0.2:8080",
  "endpoint": "/api/v2",
  "version": "2.0",
  "workspace_id": 2
}
```

#### Delete Application

```
DELETE /apps/{id}
```

### Role-Based Access Control

#### Create Workspace Role

```
POST /workspace-roles
Body: {
  "user_id": 1,
  "role": "admin",
  "workspace_id": 1
}
```

#### Get Workspace Roles

```
GET /workspace-roles
```

#### Update Workspace Role

```
PUT /workspace-roles/{id}
Body: {
  "user_id": 1,
  "role": "developer",
  "workspace_id": 1
}
```

#### Delete Workspace Role

```
DELETE /workspace-roles/{id}
```

#### Create Application Role

```
POST /app-roles
Body: {
  "user_id": 1,
  "role": "developer",
  "app_id": 1
}
```

#### Get Application Roles

```
GET /app-roles
```

#### Update Application Role

```
PUT /app-roles/{id}
Body: {
  "user_id": 1,
  "role": "admin",
  "app_id": 1
}
```

#### Delete Application Role

```
DELETE /app-roles/{id}
```

## IP Pool Management

The system includes an IP Pool management feature that automatically allocates and releases IP addresses for workspaces.

## Getting Started

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Run the server: `go run main.go`

## Testing

Run the tests using: `go test ./...`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
