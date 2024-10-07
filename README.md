# 🚀 Micro-Discover API

Welcome to the Micro-Discover API! This powerful service helps you manage users, workspaces, apps, and their respective roles. Let's dive into the different operations available!

## 👤 User Operations

### Create a User
- **POST /users**
- Description: Create a new user
- Example Request:
  ```json
  {
    "username": "john@example.com",
    "password": "securepassword123"
  }
  ```
- Example Response:
  ```json
  {
    "id": 1,
    "username": "john@example.com"
  }
  ```

### Get All Users
- **GET /users**
- Description: Retrieve all users
- Example Response:
  ```json
  [
    {
      "id": 1,
      "username": "john@example.com"
    },
    {
      "id": 2,
      "username": "jane@example.com"
    }
  ]
  ```

### Get a User
- **GET /users/{id}**
- Description: Retrieve a specific user
- Example Response:
  ```json
  {
    "id": 1,
    "username": "john@example.com"
  }
  ```

### Update a User
- **PUT /users/{id}**
- Description: Update a user's information
- Example Request:
  ```json
  {
    "username": "john.doe@example.com",
    "password": "newpassword456"
  }
  ```
- Example Response:
  ```json
  {
    "id": 1,
    "username": "john.doe@example.com"
  }
  ```

### Delete a User
- **DELETE /users/{id}**
- Description: Delete a user
- Response: 204 No Content

## 🏢 Workspace Operations

### Create a Workspace
- **POST /workspaces**
- Description: Create a new workspace
- Example Request:
  ```json
  {
    "name": "Project Alpha",
    "user_id": 1
  }
  ```
- Example Response:
  ```json
  {
    "id": 1,
    "name": "Project Alpha",
    "user_id": 1,
    "subdomain": "alpha123",
    "ips": ["10.0.0.1"]
  }
  ```

### Get All Workspaces
- **GET /workspaces**
- Description: Retrieve all workspaces
- Example Response:
  ```json
  [
    {
      "id": 1,
      "name": "Project Alpha",
      "user_id": 1,
      "subdomain": "alpha123",
      "ips": ["10.0.0.1"]
    },
    {
      "id": 2,
      "name": "Project Beta",
      "user_id": 2,
      "subdomain": "beta456",
      "ips": ["10.0.0.2"]
    }
  ]
  ```

### Get a Workspace
- **GET /workspaces/{id}**
- Description: Retrieve a specific workspace
- Example Response:
  ```json
  {
    "id": 1,
    "name": "Project Alpha",
    "user_id": 1,
    "subdomain": "alpha123",
    "ips": ["10.0.0.1"]
  }
  ```

### Update a Workspace
- **PUT /workspaces/{id}**
- Description: Update a workspace's information
- Example Request:
  ```json
  {
    "name": "Project Alpha Prime",
    "user_id": 1
  }
  ```
- Example Response:
  ```json
  {
    "id": 1,
    "name": "Project Alpha Prime",
    "user_id": 1,
    "subdomain": "alpha123",
    "ips": ["10.0.0.1"]
  }
  ```

### Delete a Workspace
- **DELETE /workspaces/{id}**
- Description: Delete a workspace
- Response: 204 No Content

## 📱 App Operations

### Create an App
- **POST /apps**
- Description: Create a new app
- Example Request:
  ```json
  {
    "name": "Awesome App",
    "description": "An awesome application",
    "git_hash": "abc123",
    "ip_port": "10.0.0.1:8080",
    "endpoint": "/api",
    "version": "1.0.0",
    "workspace_id": 1
  }
  ```
- Example Response:
  ```json
  {
    "id": 1,
    "name": "Awesome App",
    "description": "An awesome application",
    "git_hash": "abc123",
    "ip_port": "10.0.0.1:8080",
    "endpoint": "/api",
    "version": "1.0.0",
    "workspace_id": 1
  }
  ```

### Get All Apps
- **GET /apps**
- Description: Retrieve all apps
- Example Response:
  ```json
  [
    {
      "id": 1,
      "name": "Awesome App",
      "description": "An awesome application",
      "git_hash": "abc123",
      "ip_port": "10.0.0.1:8080",
      "endpoint": "/api",
      "version": "1.0.0",
      "workspace_id": 1
    },
    {
      "id": 2,
      "name": "Cool App",
      "description": "A cool application",
      "git_hash": "def456",
      "ip_port": "10.0.0.2:8080",
      "endpoint": "/api",
      "version": "2.0.0",
      "workspace_id": 2
    }
  ]
  ```

### Get an App
- **GET /apps/{id}**
- Description: Retrieve a specific app
- Example Response:
  ```json
  {
    "id": 1,
    "name": "Awesome App",
    "description": "An awesome application",
    "git_hash": "abc123",
    "ip_port": "10.0.0.1:8080",
    "endpoint": "/api",
    "version": "1.0.0",
    "workspace_id": 1
  }
  ```

### Update an App
- **PUT /apps/{id}**
- Description: Update an app's information
- Example Request:
  ```json
  {
    "name": "Super Awesome App",
    "description": "An incredibly awesome application",
    "git_hash": "abc789",
    "ip_port": "10.0.0.1:8080",
    "endpoint": "/api/v2",
    "version": "1.1.0",
    "workspace_id": 1
  }
  ```
- Example Response:
  ```json
  {
    "id": 1,
    "name": "Super Awesome App",
    "description": "An incredibly awesome application",
    "git_hash": "abc789",
    "ip_port": "10.0.0.1:8080",
    "endpoint": "/api/v2",
    "version": "1.1.0",
    "workspace_id": 1
  }
  ```

### Delete an App
- **DELETE /apps/{id}**
- Description: Delete an app
- Response: 204 No Content

## 🔑 Role Operations

### Create a Workspace Role
- **POST /workspace-roles**
- Description: Create a new workspace role
- Example Request:
  ```json
  {
    "user_id": 1,
    "role": "admin",
    "workspace_id": 1
  }
  ```
- Example Response:
  ```json
  {
    "id": 1,
    "user_id": 1,
    "role": "admin",
    "workspace_id": 1
  }
  ```

### Get All Workspace Roles
- **GET /workspace-roles**
- Description: Retrieve all workspace roles
- Example Response:
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

### Update a Workspace Role
- **PUT /workspace-roles/{id}**
- Description: Update a workspace role
- Example Request:
  ```json
  {
    "user_id": 1,
    "role": "owner",
    "workspace_id": 1
  }
  ```
- Example Response:
  ```json
  {
    "id": 1,
    "user_id": 1,
    "role": "owner",
    "workspace_id": 1
  }
  ```

### Delete a Workspace Role
- **DELETE /workspace-roles/{id}**
- Description: Delete a workspace role
- Response: 204 No Content

### Create an App Role
- **POST /app-roles**
- Description: Create a new app role
- Example Request:
  ```json
  {
    "user_id": 1,
    "role": "developer",
    "app_id": 1
  }
  ```
- Example Response:
  ```json
  {
    "id": 1,
    "user_id": 1,
    "role": "developer",
    "app_id": 1
  }
  ```

### Get All App Roles
- **GET /app-roles**
- Description: Retrieve all app roles
- Example Response:
  ```json
  [
    {
      "id": 1,
      "user_id": 1,
      "role": "developer",
      "app_id": 1
    },
    {
      "id": 2,
      "user_id": 2,
      "role": "tester",
      "app_id": 1
    }
  ]
  ```

### Update an App Role
- **PUT /app-roles/{id}**
- Description: Update an app role
- Example Request:
  ```json
  {
    "user_id": 1,
    "role": "lead developer",
    "app_id": 1
  }
  ```
- Example Response:
  ```json
  {
    "id": 1,
    "user_id": 1,
    "role": "lead developer",
    "app_id": 1
  }
  ```

### Delete an App Role
- **DELETE /app-roles/{id}**
- Description: Delete an app role
- Response: 204 No Content

## 🎉 Conclusion

That's it! You now have a complete overview of all the operations available in the Micro-Discover API. Happy coding! 🚀👨‍💻👩‍💻
