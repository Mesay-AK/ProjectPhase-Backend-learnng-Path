# Task Manager Clean Architecture

## Overview

This project is a task management system following Clean Architecture principles. It provides endpoints for user authentication, user management, and task management using the Gin framework for Go.

## Table of Contents

- [Features](#features)
- [API Endpoints](#api-endpoints)
  - [User Authentication](#user-authentication)
  - [User Management](#user-management)
  - [Task Management](#task-management)
- [Setup and Installation](#setup-and-installation)
- [Project Structure](#project-structure)


## Features

- **User Authentication**: Register and login users with JWT-based authentication.
- **User Management**: CRUD operations for user profiles, including promoting users.
- **Task Management**: CRUD operations for tasks, including creation, updating, retrieval, and deletion.

## API Endpoints

### User Authentication

- **Register User**: `POST /register`
  - **Description**: Registers a new user.
  - **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "password123",
      "name": "User Name"
    }
    ```

- **Login User**: `POST /login`
  - **Description**: Logs in a user and provides JWT tokens.
  - **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```

### User Management

- **Get All Users**: `GET /user`
  - **Description**: Retrieves a list of all users.
  
- **Get User by ID**: `GET /user/:id`
  - **Description**: Retrieves a user by their ID.

- **Promote User**: `PATCH /admin/promote/:id`
  - **Description**: Promotes a user to an admin role.
  - **Request Body**:
    ```json
    {
      "user_type": "ADMIN"
    }
    ```

- **Delete User**: `DELETE /admin/users/:id`
  - **Description**: Deletes a user by their ID.

### Task Management

- **Add Task**: `POST /admin/tasks`
  - **Description**: Creates a new task.
  - **Request Body**:
    ```json
    {
      "title": "Task Title",
      "description": "Task Description",
      "due_date": "2024-12-31T23:59:59Z"
    }
    ```

- **Get All Tasks**: `GET /user/tasks` or `GET /admin/tasks`
  - **Description**: Retrieves a list of all tasks.

- **Get Task by ID**: `GET /user/tasks/:id` or `GET /admin/tasks/:id`
  - **Description**: Retrieves a task by its ID.

- **Update Task**: `PUT /admin/tasks/:id`
  - **Description**: Updates a task's details.
  - **Request Body**:
    ```json
    {
      "title": "Updated Task Title",
      "description": "Updated Task Description",
      "due_date": "2024-12-31T23:59:59Z"
    }
    ```

- **Delete Task**: `DELETE /admin/tasks/:id`
  - **Description**: Deletes a task by its ID.

## Setup and Installation

1. **Clone the Repository**:
   ```bash
   https://github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/tree/main/Task-7

    cd task_manager_clean_architecture
    ```

2. **Install dependencies**:
   ```bash
   go mod download
    ```


3. **Set Up enviroment vairables**:
   ```bash
    SECRET_KEY=your_secret_key
    PORT=8080
    ```


4. **Run the program**:
   ```bash
   go run main.go
    ```



## Project Structure

```plaintext
.
├── Delivery
│   ├── controllers
│   ├── routers
├── Domain
│   ├── user.go
│   ├── task.go
├── Infrastructure
│   ├── auth_middleware.go
│   ├── database.go
│   ├── password_utils.go
├── Repositories
│   ├── user_repository.go
│   ├── task_repository.go
├── Usecases
│   ├── user_usecase.go
│   ├── task_usecase.go
├── main.go
└── .env

```
