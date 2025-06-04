package db

import (
	"database/sql"
	"errors"
	"fmt"

	"tidyTodo/internal/task"

	_ "modernc.org/sqlite"
)

/*
CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title VARCHAR(255) NOT NULL,
	done BOOLEAN NOT NULL DEFAULT 0
);
*/

type DB struct {
	db *sql.DB
}

func Open() (*DB, error) {
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

func (d DB) GetTasks() ([]task.Task, error) {
	rows, err := d.db.Query(`SELECT * FROM tasks;`)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var tasks []task.Task = make([]task.Task, 0)

	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			return nil, fmt.Errorf("scan tasks: %w", err)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (d *DB) InsertTask(title string) (task.Task, error) {
	if d == nil {
		return task.Task{}, errors.New("nil instance")
	}

	res, err := d.db.Exec(
		`INSERT INTO tasks (title) VALUES (?)`,
		title,
	)
	if err != nil {
		return task.Task{}, fmt.Errorf("insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return task.Task{}, fmt.Errorf("get task id: %w", err)
	}

	return task.Task{
		ID:    id,
		Title: title,
	}, nil
}

func (d *DB) UpdateTask(tsk task.Task) error {
	if d == nil {
		return errors.New("nil instance")
	}

	_, err := d.db.Exec(
		`UPDATE tasks
		SET done = ?, title = ?
		WHERE id = ?`,
		tsk.Done,
		tsk.Title,
		tsk.ID,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	return nil
}

func (d *DB) UpdateTaskStatus(id int64, done bool) error {
	if d == nil {
		return errors.New("nil instance")
	}

	_, err := d.db.Exec(
		`UPDATE tasks
		SET done = ?
		WHERE id = ?`,
		done,
		id,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}

func (d *DB) DeleteTask(id int64) error {
	if d == nil {
		return errors.New("nil instance")
	}

	_, err := d.db.Exec(
		`DELETE FROM tasks
		WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

func (d *DB) Close() error {
	if d == nil {
		return errors.New("nil instance")
	}
	return d.db.Close()
}

func (d *DB) IsClosed() bool {
	if d == nil {
		return true
	}

	err := d.db.Ping()
	if err != nil {
		return true
	}

	return false
}
