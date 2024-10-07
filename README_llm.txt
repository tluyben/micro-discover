micro-discover: Go-based microservices discovery system

Components:
1. User Service
2. Workspace Service
3. App Service
4. Role Service

Features:
- User management (CRUD)
- Workspace management (CRUD)
- App management (CRUD)
- Role management for workspaces and apps (CRUD)
- IP allocation for workspaces
- Subdomain generation for workspaces

Usage:
1. Initialize database
2. Start server
3. Use HTTP endpoints for CRUD operations

Endpoints:
- Users: /users
- Workspaces: /workspaces
- Apps: /apps
- Workspace Roles: /workspace-roles
- App Roles: /app-roles

Data Models:
- User: id, username, password
- Workspace: id, name, user_id, subdomain, ips
- App: id, name, description, git_hash, ip_port, endpoint, version, workspace_id, input_schema, output_schema
- WorkspaceRole: id, user_id, role, workspace_id
- AppRole: id, user_id, role, app_id

Authentication: Not implemented, add as needed

IP Pool: Manages IP allocation for workspaces

Subdomain Generation: Automatic for workspaces

Database: SQLite

Testing: Unit tests provided for all CRUD operations

Note: Implement proper error handling and authentication for production use.