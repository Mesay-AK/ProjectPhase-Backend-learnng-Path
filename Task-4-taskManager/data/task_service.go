package data

import (
	"errors"
	"task_manager/models"
)

var Tasks []models.Task

var ErrTaskNotFound = errors.New("error Found")

func GetTasks() ([]models.Task) {
	return Tasks
}

func GetTaskByID(id string) (models.Task, error) {
	for _, task := range Tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return models.Task{}, ErrTaskNotFound
}

func AddTask(task models.Task) error {
	Tasks = append(Tasks, task)
	return nil
}

func UpdateTask(id string, updatedTask models.Task) error {
	for i, task := range Tasks {
		if task.ID == id {
			if updatedTask.Title != "" {
				Tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				Tasks[i].Description = updatedTask.Description
			}
			if !updatedTask.DueDate.IsZero()  {
				Tasks[i].DueDate = updatedTask.DueDate
			}
			if updatedTask.Status != "" {
				Tasks[i].Status = updatedTask.Status
			}
			return nil
		}
	}
	return ErrTaskNotFound
}


func DeleteTask(id string) error {
	for i, task := range Tasks {
		if task.ID == id {
			Tasks = append(Tasks[:i], Tasks[i+1:]...)
			return nil
		}
	}
	return ErrTaskNotFound
}
