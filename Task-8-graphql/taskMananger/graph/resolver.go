package graph

import "taskManager/graph/model"

// Resolver struct will store tasks in memory
type Resolver struct {
	Tasks []*model.Task
}
