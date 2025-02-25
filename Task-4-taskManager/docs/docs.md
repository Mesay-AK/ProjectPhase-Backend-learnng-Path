# Task Manager API Documentation

## Add Task
This endpoint adds a new task.

### Request
The request should be sent via an HTTP POST method to http://localhost:8080/tasks.
The request payload should have the following parameters in the raw request body type:
  - id (string, optional): The ID of the task.
  - title (string, optional): The title of the task.
  - description (string, optional): The description of the task.
  - status (string, optional): The status of the task.

### Response 
Upon a successful execution, the endpoint returns a status code of 201 and a JSON response with a message "Task created".
## Get Tasks
This endpoint retrieves a list of tasks.

### Request
The request should be sent via an HTTP GET method to http://localhost:8080/tasks.

### Response
Upon a successful execution, the endpoint returns a status code of 200 and a JSON response with an array of tasks. Each task object in the array contains the following properties:
  - id (string): The ID of the task.
  - title (string): The title of the task.
  - description (string): The description of the task.
  - due_date (string): The due date of the task.
  - status (string): The status of the task.

### Possible Errors
  - 500 Internal Server Error: If there is an issue on the server side.

## Get Task by ID
This endpoint retrieves a task with a specific id.

### Request 
The request should be sent via an HTTP GET method to http://localhost:8080/tasks/:id.

### Response 
Upon a successful execution, the endpoint returns a status code of 200 and a JSON response with a task of that ID.
The task object contains the following properties:
  - id (string): The ID of the task.
  - title (string): The title of the task.
  - description (string): The description of the task.
  - due_date (string): The due date of the task.
  - status (string): The status of the task.

### Possible Errors
  - 404 Not Found: If the task with the specified ID is not found.
  - 500 Internal Server Error: If there is an issue on the server side.


### Possible Errors
  - 400 Bad Request: If the request body is invalid.
  - 500 Internal Server Error: If there is an issue on the server side.

## Remove Task 
This endpoint removes a task with a specific id.

### Request 
The request should be sent via an HTTP DELETE method to http://localhost:8080/tasks/:id.

### Response 
Upon a successful execution, the endpoint returns a status code of 200 and a JSON response with a message "Task removed".

### Possible Errors
  - 404 Not Found: If the task with the specified ID is not found.
  - 500 Internal Server Error: If there is an issue on the server side.

## Update Task
This endpoint updates a task of specific id.

### Request 
The request should be sent via an HTTP PUT method to http://localhost:8080/tasks/:id.
The request payload should have one or more of the following parameters in the raw request body type:
  - id (string, optional): The ID of the task.
  - title (string, optional): The title of the task.
  - description (string, optional): The description of the task.
  - status (string, optional): The status of the task.

### Response 
Upon a successful execution, the endpoint returns a status code of 200 and a JSON response with a message "Task updated".

### Possible Errors
  - 400 Bad Request: If the request body is invalid.
  - 404 Not Found: If the task with the specified ID is not found.
  - 500 Internal Server Error: If there is an issue on the server side.
