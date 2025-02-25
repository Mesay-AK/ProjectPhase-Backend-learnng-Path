package usecase

import (
	"context"
	repository "task_manager_clean/Repositories"
	domain "task_manager_clean/Domain"
	"github.com/go-playground/validator/v10"
)

type TaskUseCase struct {
	TaskRepo *repository.TaskRepository
	Validate *validator.Validate
}

func NewTaskUseCase(tr *repository.TaskRepository, v *validator.Validate) *TaskUseCase {
	return &TaskUseCase{
		TaskRepo: tr,
		Validate: v,
	}
}


func (uCase *TaskUseCase) GetTasks(ctx context.Context) ([]domain.Task, error) {
	return uCase.TaskRepo.GetTasks(ctx) 
}

func (uCase *TaskUseCase) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	task, err := uCase.TaskRepo.GetTaskByID(ctx, id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			return nil, domain.ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}

func (uCase *TaskUseCase) AddTask(ctx context.Context, task *domain.Task) error {

	if err := uCase.Validate.Struct(task); err != nil {
		return err
	}
	return uCase.TaskRepo.AddTask(ctx, task)
}

func (uCase *TaskUseCase) UpdateTask(ctx context.Context, id string, updatedTask *domain.Task) (*domain.Task, error) {
	
	if err := uCase.Validate.Struct(updatedTask); err != nil {
		return nil, err
	}
	return uCase.TaskRepo.UpdateTask(ctx, id, updatedTask)
}

func (uCase *TaskUseCase) DeleteTask(ctx context.Context, id string) error {
	return uCase.TaskRepo.DeleteTask(ctx, id)
}
