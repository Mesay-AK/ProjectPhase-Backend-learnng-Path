package graph

import (
	"context"
	"taskManager/internal/models"
	"taskManager/graph/model" 
	"time"
	"strconv"
)

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }



func convertToGraphModel(task models.Task) *model.Task {

	createdAtStr := task.CreatedAt.Format(time.RFC3339)
	description := &task.Description
	if task.Description == "" {
		description = nil
	}

	return &model.Task{
		ID:          strconv.FormatUint(uint64(task.ID), 10), 
		Title:       task.Title,
		Description: description, 
		Completed:   task.Completed,
		CreatedAt:   createdAtStr,
	}
}


func (r *mutationResolver) AddTask(ctx context.Context, title string, description *string) (*model.Task, error) {
	task := models.Task{
		Title:       title,
		Description: *description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	if err := r.DB.Create(&task).Error; err != nil {
		return nil, err
	}

	return convertToGraphModel(task), nil
}

func (r *mutationResolver) UpdateTask(ctx context.Context, id string, title *string, description *string, completed *bool) (*model.Task, error) {
	var task models.Task
	if err := r.DB.First(&task, id).Error; err != nil {
		return nil, err
	}

	if title != nil {
		task.Title = *title
	}
	if description != nil {
		task.Description = *description
	}
	if completed != nil {
		task.Completed = *completed
	}

	if err := r.DB.Save(&task).Error; err != nil {
		return nil, err
	}

	// Convert models.Task to model.Task and return
	return convertToGraphModel(task), nil
}

// DeleteTask is the resolver for the deleteTask field.
func (r *mutationResolver) DeleteTask(ctx context.Context, id string) (bool, error) {
	if err := r.DB.Delete(&models.Task{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

// GetTasks is the resolver for the getTasks field.
func (r *queryResolver) GetTasks(ctx context.Context) ([]*model.Task, error) {
	var tasks []*models.Task
	if err := r.DB.Find(&tasks).Error; err != nil {
		return nil, err
	}

	// Convert []*models.Task to []*model.Task
	var result []*model.Task
	for _, task := range tasks {
		result = append(result, convertToGraphModel(*task))
	}
	return result, nil
}

// GetTask is the resolver for the getTask field.
func (r *queryResolver) GetTask(ctx context.Context, id string) (*model.Task, error) {
	var task models.Task
	if err := r.DB.First(&task, id).Error; err != nil {
		return nil, err
	}
	return convertToGraphModel(task), nil
}



