package common

import (
	"errors"
	"fmt"
	"tidyTodo/internal/db"
	"tidyTodo/internal/task"
)

var DB *db.DB

func UpdateTaskStatus(tsk *task.Task) error {
	if DB.IsClosed() {
		return errors.New("db closed")
	}

	if err := tsk.Check(); err != nil {
		return fmt.Errorf("update task status: check: %w", err)
	}
	if err := DB.UpdateTaskStatus(tsk.ID, tsk.Done); err != nil {
		return fmt.Errorf("update task status: db: %w", err)
	}
	return nil
}
