package main

import (
	"fmt"
	"io"
	"log"
	"task-manager/graphql"
	"encoding/json"
)

// Task struct to represent a task
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// GetTasks fetches all tasks from Hasura
func GetTasks() {
	resp, err := graphql.ExecuteQuery(graphql.GetTasksQuery, nil)
	if err != nil {
		log.Fatal("Failed to fetch tasks:", err)
	}

	defer resp.Body.Close()

	var result struct {
		Data struct {
			Tasks []Task `json:"tasks"`
		} `json:"data"`
	}

	body, _ := io.ReadAll(resp.Body) // Updated to io.ReadAll
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal("Failed to decode response:", err)
	}

	fmt.Println("Tasks:", result.Data.Tasks)
}

// AddTask adds a new task to Hasura
func AddTask(title, description string) {
	variables := map[string]interface{}{
		"title":       title,
		"description": description,
	}

	resp, err := graphql.ExecuteQuery(graphql.AddTaskMutation, variables)
	if err != nil {
		log.Fatal("Failed to add task:", err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body) // Updated to io.ReadAll
	fmt.Println("Added Task:", string(body))
}

// UpdateTask updates the status of an existing task
func UpdateTask(id, status string) {
	variables := map[string]interface{}{
		"id":     id,
		"status": status,
	}

	resp, err := graphql.ExecuteQuery(graphql.UpdateTaskMutation, variables)
	if err != nil {
		log.Fatal("Failed to update task:", err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body) // Updated to io.ReadAll
	fmt.Println("Updated Task:", string(body))
}

// DeleteTask deletes a task by its ID
func DeleteTask(id string) {
	variables := map[string]interface{}{
		"id": id,
	}

	resp, err := graphql.ExecuteQuery(graphql.DeleteTaskMutation, variables)
	if err != nil {
		log.Fatal("Failed to delete task:", err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body) // Updated to io.ReadAll
	fmt.Println("Deleted Task:", string(body))
}

func main() {
	// Fetch tasks
	fmt.Println("Fetching tasks...")
	GetTasks()

	// Add a new task
	fmt.Println("Adding a new task...")
	AddTask("Learn GraphQL", "Study Hasura and Golang integration")

	// Update an existing task (You can replace the ID and status)
	fmt.Println("Updating a task...")
	UpdateTask("1d973e6a-f9f8-4f34-bda7-792ace493db3", "completed")

	// Delete a task (You can replace the ID)
	fmt.Println("Deleting a task...")
	DeleteTask("1d973e6a-f9f8-4f34-bda7-792ace493db3")
}
