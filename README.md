# 🚀 Micro-Discover

Welcome to Micro-Discover, your ultimate solution for managing microservices! 🎉

## 📋 Features

- User Management 👤
- Workspace Management 🏢
- App Management 📱
- Role-based Access Control 🔐

## 🗄️ Database Schema

### Users Table
- id: INTEGER (Primary Key)
- username: TEXT
- password: TEXT

### Workspaces Table
- id: INTEGER (Primary Key)
- name: TEXT
- user_id: INTEGER (Foreign Key to Users)
- subdomain: TEXT
- ips: TEXT

### Apps Table
- id: INTEGER (Primary Key)
- name: TEXT
- description: TEXT
- git_hash: TEXT
- ip_port: TEXT
- endpoint: TEXT
- version: TEXT
- workspace_id: INTEGER (Foreign Key to Workspaces)
- input_schema: TEXT (Stringified JSON)
- output_schema: TEXT (Stringified JSON)

### Workspace Roles Table
- id: INTEGER (Primary Key)
- user_id: INTEGER (Foreign Key to Users)
- role: TEXT
- workspace_id: INTEGER (Foreign Key to Workspaces)

### App Roles Table
- id: INTEGER (Primary Key)
- user_id: INTEGER (Foreign Key to Users)
- role: TEXT
- app_id: INTEGER (Foreign Key to Apps)

## 📊 Input/Output Schemas

### Input Schema
The input schema is a stringified JSON object that describes the expected input for an app. It's stored in the `input_schema` column of the Apps table.

### Output Schema
The output schema is a stringified JSON object that describes the expected output from an app. It's stored in the `output_schema` column of the Apps table. The new output schema includes:

```json
{
  "output": [
    {
      "name": "answer",
      "type": "string"
    },
    {
      "name": "confidence",
      "type": "number"
    }
  ]
}
```

## 🚀 Getting Started

1. Clone the repository
2. Install dependencies
3. Set up the database
4. Run the server

## 🛠️ API Endpoints

- `/users`: User management
- `/workspaces`: Workspace management
- `/apps`: App management
- `/workspace-roles`: Workspace role management
- `/app-roles`: App role management

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines for more details.

## 📄 License

This project is licensed under the MIT License.

Happy coding! 💻🎉