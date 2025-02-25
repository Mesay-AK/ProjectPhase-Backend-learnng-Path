package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("Title is required")
	}
	if t.Description == "" {
		return errors.New("Description is required")
	}
	if t.Status == "" {
		return errors.New("Status is required")
	}
	if t.Due_Date.IsZero() {
		return errors.New("Due_Date is required")
	}
	return nil
}

func TestTaskValidation(t *testing.T) {
	t.Run("Test should pass with a valid task", func(t *testing.T) {
		newTask := Task{
			ID:          primitive.NewObjectID(),
			Title:       "Test Task",
			Description: "This is a test task",
			Status:      "Started",
			Due_Date:    time.Now(),
		}

		assert.NotEmpty(t, newTask.ID, "ID should not be empty")
		assert.Equal(t, "Test Task", newTask.Title, "Title should be Test Task")
		assert.Equal(t, "This is a test task", newTask.Description, "Description should match")
		assert.Equal(t, TaskStatus("Started"), newTask.Status, "Status should be started")
		assert.NotEmpty(t, newTask.Due_Date, "due date should not be empty")
	})

	t.Run("Test should fail with a missing field", func(t *testing.T) {
		newTask := Task{
			ID:          primitive.NewObjectID(),
			Title:       "Test Task",
			Description: "This is a test task",
			Status:      "Started",
		}

		err := newTask.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Due_Date")
	})

}
