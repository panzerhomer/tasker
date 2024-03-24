# Tasker

Project management tool
REST API

# Installation
```
docker-comose up -d
go run main.go
```

## Authentication
- **POST /login**:
  - Allows users to authenticate in the system.
- **POST /signup**:
  - Registers new users in the system.

## Users
- **GET /hello**:
  - Returns a welcome message or information about the current user.
- **GET /logout**:
  - Ends the user session.

## Projects
- **POST /projects**:
  - Creates a new project.
- **GET /projects**:
  - Returns a list of all projects.
- **GET /projects/{name}/users**:
  - Retrieves a list of users for a specific project.
- **PUT /projects/{name}**:
  - Adds a user to the project.
- **DELETE /projects/{name}**:
  - Removes a user from the project.

## Tasks
- **POST /projects/{name}/tasks**:
  - Creates a new task within a project.
- **GET /projects/{name}/tasks**:
  - Returns a list of tasks for the project.

## Search
- **GET /projects/search**:
  - Search project by name usign page


 simple rest api 
