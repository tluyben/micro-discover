# Discovery Service for Sandboxed Applications

## Overview

The Discovery Service is a Go-based application designed to manage users, workspaces, and applications in a sandboxed environment. It provides a RESTful API for creating, reading, updating, and deleting (CRUD) users, workspaces, and applications. The service uses SQLite for data storage and implements IP allocation and subdomain generation for workspaces.

## Features

- User management with secure password hashing
- Workspace management with automatic IP allocation and subdomain generation
- Application management within workspaces
- RESTful API for all operations
- SQLite database for persistent storage
- IP pool management for efficient IP allocation
- Concurrent-safe operations

## Prerequisites

- Go 1.16 or higher
- SQLite3

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/yourusername/discovery-service.git
   cd discovery-service
   ```

2. Install dependencies:

   ```
   go get github.com/gorilla/mux
   go get github.com/mattn/go-sqlite3
   go get golang.org/x/crypto/bcrypt
   ```

3. Build the application:
   ```
   go build -o discovery-service
   ```

## Usage

1. Start the service:

   ```
   ./discovery-service
   ```

   The service will start on port 8080 and create a SQLite database file named `discovery.db` in the same directory.

2. Use the API endpoints to interact with the service.

## API Endpoints

### Users

- `POST /users`: Create a new user
- `GET /users`: List all users
- `GET /users/{id}`: Get a specific user
- `PUT /users/{id}`: Update a user
- `DELETE /users/{id}`: Delete a user

### Workspaces

- `POST /workspaces`: Create a new workspace
- `GET /workspaces`: List all workspaces
- `GET /workspaces/{id}`: Get a specific workspace
- `PUT /workspaces/{id}`: Update a workspace
- `DELETE /workspaces/{id}`: Delete a workspace

### Applications

- `POST /apps`: Create a new application
- `GET /apps`: List all applications
- `GET /apps/{id}`: Get a specific application
- `PUT /apps/{id}`: Update an application
- `DELETE /apps/{id}`: Delete an application

## Data Models

### User

```json
{
  "id": 1,
  "username": "user@example.com"
}
```

### Workspace

```json
{
  "id": 1,
  "name": "example_workspace",
  "user_id": 1,
  "subdomain": "abcd1234",
  "ips": ["10.0.0.1"]
}
```

### Application

```json
{
  "id": 1,
  "name": "example_app",
  "description": "An example application",
  "git_hash": "abc123def456",
  "ip_port": "10.0.0.1:8080",
  "endpoint": "/api",
  "version": "1.0.0",
  "workspace_id": 1
}
```

## Example Requests

### Create a User

```bash
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"username": "user@example.com", "password": "securepassword"}'
```

Response:
```json
{
  "id": 1,
  "username": "user@example.com"
}
```

### Create a Workspace

```bash
curl -X POST http://localhost:8080/workspaces -H "Content-Type: application/json" -d '{"name": "myworkspace", "user_id": 1}'
```

Response:
```json
{
  "id": 1,
  "name": "myworkspace",
  "user_id": 1,
  "subdomain": "abcd1234",
  "ips": ["10.0.0.1"]
}
```

### Create an Application

```bash
curl -X POST http://localhost:8080/apps -H "Content-Type: application/json" -d '{"name": "myapp", "description": "My first app", "git_hash": "abc123", "ip_port": "10.0.0.1:8080", "endpoint": "/api", "version": "1.0", "workspace_id": 1}'
```

Response:
```json
{
  "id": 1,
  "name": "myapp",
  "description": "My first app",
  "git_hash": "abc123",
  "ip_port": "10.0.0.1:8080",
  "endpoint": "/api",
  "version": "1.0",
  "workspace_id": 1
}
```

### Get All Users

```bash
curl http://localhost:8080/users
```

Response:
```json
[
  {
    "id": 1,
    "username": "user@example.com"
  }
]
```

### Get All Workspaces

```bash
curl http://localhost:8080/workspaces
```

Response:
```json
[
  {
    "id": 1,
    "name": "myworkspace",
    "user_id": 1,
    "subdomain": "abcd1234",
    "ips": ["10.0.0.1"]
  }
]
```

### Get All Applications

```bash
curl http://localhost:8080/apps
```

Response:
```json
[
  {
    "id": 1,
    "name": "myapp",
    "description": "My first app",
    "git_hash": "abc123",
    "ip_port": "10.0.0.1:8080",
    "endpoint": "/api",
    "version": "1.0",
    "workspace_id": 1
  }
]
```

## Security Considerations

- The service uses bcrypt for password hashing.
- IP allocation is managed to prevent conflicts.
- Subdomains are generated randomly and checked for uniqueness.
- The service does not implement authentication or authorization for API endpoints. In a production environment, you should add appropriate security measures.

## Limitations and Future Improvements

- The service uses SQLite, which may not be suitable for high-concurrency scenarios. Consider using a more robust database for production use.
- Error handling could be improved for better debugging and user feedback.
- Adding logging would be beneficial for monitoring and troubleshooting.
- Implementing rate limiting would help prevent abuse of the API.
- Adding metrics and health check endpoints would improve observability.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
