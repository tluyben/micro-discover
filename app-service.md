# App Service 🚀

The App Service is responsible for managing applications within the micro-discover ecosystem.

## Endpoints 🛣️

### 1. Create App 📝

**POST** `/apps`

Creates a new application.

**Request Body:**
```json
{
  "name": "MyApp",
  "description": "This is my awesome app",
  "git_hash": "abc123",
  "ip_port": "10.0.0.1:8080",
  "endpoint": "/api/v1",
  "version": "1.0.0",
  "workspace_id": 1,
  "input_schema": {
    "type": "object",
    "properties": {
      "name": {"type": "string"},
      "age": {"type": "integer"}
    }
  },
  "output_schema": {
    "type": "object",
    "properties": {
      "greeting": {"type": "string"}
    }
  }
}
```

**Response:**
```json
{
  "id": 1,
  "name": "MyApp",
  "description": "This is my awesome app",
  "git_hash": "abc123",
  "ip_port": "10.0.0.1:8080",
  "endpoint": "/api/v1",
  "version": "1.0.0",
  "workspace_id": 1,
  "input_schema": {...},
  "output_schema": {...}
}
```

### 2. Get App 🔍

**GET** `/apps/{id}`

Retrieves details of a specific application.

**Response:**
```json
{
  "id": 1,
  "name": "MyApp",
  "description": "This is my awesome app",
  "git_hash": "abc123",
  "ip_port": "10.0.0.1:8080",
  "endpoint": "/api/v1",
  "version": "1.0.0",
  "workspace_id": 1,
  "input_schema": {...},
  "output_schema": {...}
}
```

### 3. Update App 🔄

**PUT** `/apps/{id}`

Updates an existing application.

**Request Body:**
```json
{
  "name": "UpdatedApp",
  "description": "This is my updated app",
  "git_hash": "def456",
  "ip_port": "10.0.0.2:8080",
  "endpoint": "/api/v2",
  "version": "2.0.0",
  "workspace_id": 1,
  "input_schema": {...},
  "output_schema": {...}
}
```

**Response:**
```json
{
  "id": 1,
  "name": "UpdatedApp",
  "description": "This is my updated app",
  "git_hash": "def456",
  "ip_port": "10.0.0.2:8080",
  "endpoint": "/api/v2",
  "version": "2.0.0",
  "workspace_id": 1,
  "input_schema": {...},
  "output_schema": {...}
}
```

### 4. Delete App 🗑️

**DELETE** `/apps/{id}`

Deletes an application.

**Response:**
Status: 204 No Content

### 5. List Apps 📋

**GET** `/apps`

Retrieves a list of all applications.

**Response:**
```json
[
  {
    "id": 1,
    "name": "App1",
    "description": "Description 1",
    "git_hash": "abc123",
    "ip_port": "10.0.0.1:8080",
    "endpoint": "/api/v1",
    "version": "1.0.0",
    "workspace_id": 1,
    "input_schema": {...},
    "output_schema": {...}
  },
  {
    "id": 2,
    "name": "App2",
    "description": "Description 2",
    "git_hash": "def456",
    "ip_port": "10.0.0.2:8080",
    "endpoint": "/api/v2",
    "version": "2.0.0",
    "workspace_id": 2,
    "input_schema": {...},
    "output_schema": {...}
  }
]
```

## Fields 📊

- `id`: Unique identifier for the app (integer)
- `name`: Name of the app (string)
- `description`: Description of the app (string)
- `git_hash`: Git hash of the app's code (string)
- `ip_port`: IP and port where the app is running (string)
- `endpoint`: API endpoint of the app (string)
- `version`: Version of the app (string)
- `workspace_id`: ID of the workspace the app belongs to (integer)
- `input_schema`: JSON schema defining the input structure (object)
- `output_schema`: JSON schema defining the output structure (object)

## Examples 💡

### Creating an App

```bash
curl -X POST http://localhost:8080/apps \
-H "Content-Type: application/json" \
-d '{
  "name": "WeatherApp",
  "description": "Provides weather forecasts",
  "git_hash": "abc123def456",
  "ip_port": "10.0.0.1:8080",
  "endpoint": "/api/weather",
  "version": "1.0.0",
  "workspace_id": 1,
  "input_schema": {
    "type": "object",
    "properties": {
      "city": {"type": "string"},
      "country": {"type": "string"}
    }
  },
  "output_schema": {
    "type": "object",
    "properties": {
      "temperature": {"type": "number"},
      "conditions": {"type": "string"}
    }
  }
}'
```

### Updating an App

```bash
curl -X PUT http://localhost:8080/apps/1 \
-H "Content-Type: application/json" \
-d '{
  "name": "WeatherApp Pro",
  "description": "Advanced weather forecasts",
  "git_hash": "ghi789jkl012",
  "ip_port": "10.0.0.2:8080",
  "endpoint": "/api/weather/pro",
  "version": "2.0.0",
  "workspace_id": 1,
  "input_schema": {
    "type": "object",
    "properties": {
      "city": {"type": "string"},
      "country": {"type": "string"},
      "days": {"type": "integer"}
    }
  },
  "output_schema": {
    "type": "object",
    "properties": {
      "forecast": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "date": {"type": "string"},
            "temperature": {"type": "number"},
            "conditions": {"type": "string"}
          }
        }
      }
    }
  }
}'
```

These examples demonstrate how to interact with the App Service API to create and update applications. 🎉
