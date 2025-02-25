## Project structure

```
Task_manager/
├── Bootstrap/
│   ├── app.go
│   ├── database.go
│   └── env.go
├── Delivery/
│   ├── controllers/
│   │   ├── login_controller.go
│   │   ├── login_controller_test.go
│   │   ├── promote_controller.go
│   │   ├── promote_controller_test.go
│   │   ├── signup_controller.go
│   │   ├── signup_controller_test.go
│   │   ├── tasks_controller_test.go
│   │   └── tasks_controller.go
│   ├── routers/
│   │   ├── login_route.go
│   │   ├── promote_route.go
│   │   ├── route.go
│   │   ├── signup_route.go
│   │   ├── task_admin_route.go
│   │   └── task_user_route.go
│   └── main.go
├── docs/
│   └── documentation.md
├── Domain/
│   ├── jwt_custom.go
│   ├── login.go
│   ├── signup.go
│   ├── task.go
│   └── user.go
├── Infrastructure/
│   ├── middleware/
│   │   ├── admin_middleware.go
│   │   └── auth_middleware.go
│   ├── jwt_service.go
│   └── password_service.go
├── Repositories/
│   ├── task_repository_test.go
│   ├── task_repository.go
│   ├── user_repository_test.go
│   └── user_repository.go
├── Usecases/
│   ├── login_usecase_test.go
│   ├── login_usecase.go
│   ├── promote_usecase_test.go
│   ├── promote_usecase.go
│   ├── signup_usecase_test.go
│   ├── signup_usecase.go
│   ├── task_usecase_test.go
│   └── task_usecase.go
├── go.mod
└── go.sum


```

## Overview

This project is an API for a task manager application that has all the needed endpoints for creating tasks, getting all tasks, updating tasks, deleting tasks, and getting a single task by ID. It uses MongoDB for data storage and supports two roles: admin and user. Additionally, it uses JWT for authentication and authorization. It has a role based authentication and authorization and it uses clean architecture patterns.

## Setup and Installation

### Installation

1. Clone the repository:

```sh
git clone https://github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task-8.git
```
```sh
cd Task8-testing
```
2. Fetch the dependencies:

```sh
go get ./...

```

3. Initialize the project:

```sh
go mod tidy

```
4. Before running the project monog db installation is needed since this api works with a local database

- For linux installation : https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-ubuntu/
- For windows installation : https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-windows/
- For macOS installation : https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-os-x/

## Running the project

```sh
go run Delivery/main.go

```

## Running tests

- To run the tests either you can run all tests at once.

```sh
go tets ./... -v

```

- OR run particular tests one by one
```sh
go tets ./Delivery/controllers -v
go tets ./Repositories -v
go tets ./Usecases -v
go tets ./Domain -v
```


- After Monog db installation and runnig the app a mongo db database called taskManager with a collection of tasks in it will be created.

There are Five Api End Points for this backend
- Post http://localhost:4000/user/register with a user request body
- Post http://localhost:4000/user/login with a user request body
- Post http://localhost:4000/user/promote/:id - needed admin token
- Get http://localhost:4000/tasks
- Get http://localhost:4000/tasks/2
- Post http://localhost:4000/tasks with a task request body - needed admin token
- Put http://localhost:4000/tasks/2 with a task request body - needed admin token
- Delete http://localhost:4000/tasks/2 - needed admin token

Detailed Api Documentation - https://documenter.getpostman.com/view/29564648/2sA3s3GqWn