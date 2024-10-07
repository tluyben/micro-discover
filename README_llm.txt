# Micro-Discover API Documentation 🚀

This document provides a comprehensive guide to the Micro-Discover REST API endpoints.

## Users 👤

### Create User
POST /users
Body: {"username": "string", "password": "string"}
Response: {"id": int, "username": "string"}
Creates a new user with the given username (email) and password.

### Get Users
GET /users
Response: [{"id": int, "username": "string"}]
Returns a list of all users.

### Get User
GET /users/{id}
Response: {"id": int, "username": "string"}
Returns details of a specific user.

### Update User
PUT /users/{id}
Body: {"username": "string", "password": "string"}
Response: {"id": int, "username": "string"}
Updates the details of a specific user.

### Delete User
DELETE /users/{id}
Deletes a specific user.

## Workspaces 🏢

### Create Workspace
POST /workspaces
Body: {"name": "string", "user_id": int}
Response: {"id": int, "name": "string", "user_id": int, "subdomain": "string", "ips": ["string"]}
Creates a new workspace for a user.

### Get Workspaces
GET /workspaces
Response: [{"id": int, "name": "string", "user_id": int, "subdomain": "string", "ips": ["string"]}]
Returns a list of all workspaces.

### Get Workspace
GET /workspaces/{id}
Response: {"id": int, "name": "string", "user_id": int, "subdomain": "string", "ips": ["string"]}
Returns details of a specific workspace.

### Update Workspace
PUT /workspaces/{id}
Body: {"name": "string", "user_id": int}
Response: {"id": int, "name": "string", "user_id": int, "subdomain": "string", "ips": ["string"]}
Updates the details of a specific workspace.

### Delete Workspace
DELETE /workspaces/{id}
Deletes a specific workspace and releases its IPs.

## Apps 📱

### Create App
POST /apps
Body: {"name": "string", "description": "string", "git_hash": "string", "ip_port": "string", "endpoint": "string", "version": "string", "workspace_id": int, "input_schema": "string", "output_schema": "string"}
Response: App object
Creates a new app in a workspace.

### Get Apps
GET /apps
Response: [App objects]
Returns a list of all apps.

### Get App
GET /apps/{id}
Response: App object
Returns details of a specific app.

### Update App
PUT /apps/{id}
Body: App object (same as create)
Response: Updated App object
Updates the details of a specific app.

### Delete App
DELETE /apps/{id}
Deletes a specific app.

## Workspace Roles 🔑

### Create Workspace Role
POST /workspace-roles
Body: {"user_id": int, "role": "string", "workspace_id": int}
Response: {"id": int, "user_id": int, "role": "string", "workspace_id": int}
Assigns a role to a user for a specific workspace.

### Get Workspace Roles
GET /workspace-roles
Response: [{"id": int, "user_id": int, "role": "string", "workspace_id": int}]
Returns a list of all workspace roles.

### Update Workspace Role
PUT /workspace-roles/{id}
Body: {"user_id": int, "role": "string", "workspace_id": int}
Response: {"id": int, "user_id": int, "role": "string", "workspace_id": int}
Updates a specific workspace role.

### Delete Workspace Role
DELETE /workspace-roles/{id}
Removes a specific workspace role.

## App Roles 🔐

### Create App Role
POST /app-roles
Body: {"user_id": int, "role": "string", "app_id": int}
Response: {"id": int, "user_id": int, "role": "string", "app_id": int}
Assigns a role to a user for a specific app.

### Get App Roles
GET /app-roles
Response: [{"id": int, "user_id": int, "role": "string", "app_id": int}]
Returns a list of all app roles.

### Update App Role
PUT /app-roles/{id}
Body: {"user_id": int, "role": "string", "app_id": int}
Response: {"id": int, "user_id": int, "role": "string", "app_id": int}
Updates a specific app role.

### Delete App Role
DELETE /app-roles/{id}
Removes a specific app role.

This API allows for comprehensive management of users, workspaces, apps, and roles within the Micro-Discover system. Each endpoint is designed to perform specific CRUD operations on the respective entities, providing a flexible and powerful interface for interacting with the system.