package usecases

import (
	"context"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
)

type taskUsecase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

// NewTaskUsecase creates a new instance of the TaskUsecase struct.
// It takes a taskRepository of type domain.TaskRepository and a timeout of type time.Duration as parameters.
// It returns a pointer to the TaskUsecase struct.
// The taskRepository parameter is used to interact with the task repository.
// The timeout parameter specifies the maximum duration for each operation performed by the usecase.
func NewTaskUsecase(taskRepository domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

// CreateTask creates a new task.
// It takes a context and a task as input parameters.
// The context is used for managing the execution deadline and cancellation.
// The task parameter represents the task to be created.
// It returns an error if there was a problem creating the task.
func (tu *taskUsecase) CreateTask(c context.Context, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.CreateTask(ctx, task)
}

// Fetch retrieves a list of tasks from the task repository.
// It takes a context.Context as a parameter to handle timeouts and cancellations.
// It returns a slice of domain.Task and an error.
func (tu *taskUsecase) Fetch(c context.Context) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.Fetch(ctx)
}

// GetByID retrieves a task by its ID.
// It takes a context.Context and the ID of the task as parameters.
// It returns the retrieved task and an error, if any.
func (tu *taskUsecase) GetByID(c context.Context, id string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.GetByID(ctx, id)
}

// UpdateTask updates a task with the given task update data and ID.
// It takes a context.Context, a *domain.TaskUpdate, and a string as parameters.
// The context.Context is used for managing the execution context and timeout.
// The *domain.TaskUpdate contains the updated task data.
// The string represents the ID of the task to be updated.
// It returns a domain.Task and an error.
// The domain.Task represents the updated task.
// The error indicates if any error occurred during the update process.
func (tu *taskUsecase) UpdateTask(c context.Context, task *domain.TaskUpdate, id string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.UpdateTask(ctx, task, id)
}

// DeleteTask deletes a task with the specified ID.
// It takes a context.Context and the ID of the task to be deleted as parameters.
// It returns an error if the task deletion fails.
func (tu *taskUsecase) DeleteTask(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.DeleteTask(ctx, id)
}
