package models

import (
	"time"
)

type Task struct {
	ID          string    `bson:"_id" json:"id"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	DueDate     time.Time `bson:"due_date" json:"due_date"`
	Status      string    `bson:"status" json:"status"`
}