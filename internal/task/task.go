package task

import (
	"errors"
)

type Task struct {
	ID    int64
	Title string
	Done  bool
}

func (t *Task) Check() error {
	if t == nil || t.ID == 0 {
		return errors.New("nil instance")
	}
	t.Done = !t.Done
	return nil
}
