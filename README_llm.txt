micro-discover API:

Users:
POST /users - Create user (JSON: {username, password})
GET /users - List users
GET /users/{id} - Get user
PUT /users/{id} - Update user (JSON: {username, password})
DELETE /users/{id} - Delete user

Workspaces:
POST /workspaces - Create workspace (JSON: {name, user_id})
GET /workspaces - List workspaces
GET /workspaces/{id} - Get workspace
PUT /workspaces/{id} - Update workspace (JSON: {name, user_id})
DELETE /workspaces/{id} - Delete workspace

Apps:
POST /apps - Create app (JSON: {name, description, git_hash, ip_port, endpoint, version, workspace_id})
GET /apps - List apps
GET /apps/{id} - Get app
PUT /apps/{id} - Update app (JSON: same as POST)
DELETE /apps/{id} - Delete app

Workspace Roles:
POST /workspace-roles - Create workspace role (JSON: {user_id, role, workspace_id})
GET /workspace-roles - List workspace roles
PUT /workspace-roles/{id} - Update workspace role (JSON: same as POST)
DELETE /workspace-roles/{id} - Delete workspace role

App Roles:
POST /app-roles - Create app role (JSON: {user_id, role, app_id})
GET /app-roles - List app roles
PUT /app-roles/{id} - Update app role (JSON: same as POST)
DELETE /app-roles/{id} - Delete app role

All POST requests return 201 Created on success.
All GET requests return 200 OK on success.
All PUT requests return 200 OK on success.
All DELETE requests return 204 No Content on success.

Error responses: 400 Bad Request, 404 Not Found, 500 Internal Server Error