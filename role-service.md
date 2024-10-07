# 🎭 Role Service

The Role Service manages user roles within the micro-discover ecosystem. Roles are assigned per app or per workspace and are associated with specific users.

## 📊 Role Types

- **Workspace Roles**: Roles that apply to an entire workspace
- **App Roles**: Roles that apply to a specific app within a workspace

## 🔐 Available Roles

### Workspace Roles
- Admin: Full control over the workspace
- Member: Basic access to the workspace

### App Roles
- Developer: Can modify and deploy the app
- User: Can use the app

## 🛠️ Endpoints

### Workspace Roles

#### Create Workspace Role
- **POST** `/workspace-roles`
- Body:
  ```json
  {
    "user_id": 1,
    "role": "admin",
    "workspace_id": 1
  }
  ```

#### Get Workspace Roles
- **GET** `/workspace-roles`

#### Update Workspace Role
- **PUT** `/workspace-roles/{id}`
- Body:
  ```json
  {
    "user_id": 1,
    "role": "member",
    "workspace_id": 1
  }
  ```

#### Delete Workspace Role
- **DELETE** `/workspace-roles/{id}`

### App Roles

#### Create App Role
- **POST** `/app-roles`
- Body:
  ```json
  {
    "user_id": 1,
    "role": "developer",
    "app_id": 1
  }
  ```

#### Get App Roles
- **GET** `/app-roles`

#### Update App Role
- **PUT** `/app-roles/{id}`
- Body:
  ```json
  {
    "user_id": 1,
    "role": "user",
    "app_id": 1
  }
  ```

#### Delete App Role
- **DELETE** `/app-roles/{id}`

## 🔗 Integration

The Role Service integrates closely with the User Service and Workspace Service to ensure proper access control and permissions management across the micro-discover platform.

## 🚀 Future Enhancements

- Role inheritance
- Custom role definitions
- Time-based role assignments

Remember to always use proper authentication and authorization when accessing these endpoints to maintain the security of your micro-discover deployment! 🔒👨‍💻👩‍💻