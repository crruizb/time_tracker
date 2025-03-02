package data

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

type TasksPostgres struct {
	DB *sql.DB
}

func NewTasksPostgres(db *sql.DB) *TasksPostgres {
	return &TasksPostgres{
		DB: db,
	}
}

type Task struct {
	Id          string     `json:"id" db:"id"`
	ProjectId   string     `json:"projectId" db:"project_id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Username    *string    `json:"username" db:"username"`
	StartedAt   *time.Time `json:"startedAt" db:"started_at"`
	FinishedAt  *time.Time `json:"finishedAt" db:"finished_at"`
}

var (
	ErrTaskAlreadyStarted  = errors.New("task already started")
	ErrTaskNotStarted      = errors.New("task is not started")
	ErrTaskAlreadyFinished = errors.New("task already finished")
	ErrProjectNotFound     = errors.New("project not found")
)

func (m TasksPostgres) CreateTask(projectId, name, description string) (*Task, error) {
	query := `
		INSERT INTO tasks (project_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	args := []any{projectId, name, description}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	task := &Task{
		ProjectId:   projectId,
		Name:        name,
		Description: description,
	}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&task.Id)
	if err != nil {
		slog.Error("error creating task", "err", err.Error())
		if err.Error() == `pq: insert or update on table "tasks" violates foreign key constraint "fk_project_id"` {
			return nil, ErrProjectNotFound
		}

		return nil, err
	}

	return task, nil
}

func (m TasksPostgres) StartTask(taskId, userId string) error {
	task, err := m.getTask(taskId, userId)
	if err != nil {
		return err
	}
	if task != nil {
		return ErrTaskAlreadyStarted
	}

	query := `
		INSERT INTO tasks_users (task_id, user_id, started_at)
		VALUES ($1, $2, $3);
	`

	args := []any{taskId, userId, time.Now().UTC().Format(time.RFC3339)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		slog.Error("error starting task", "err", err.Error())
		return err
	}

	return nil
}

func (m TasksPostgres) StopTask(taskId, userId string) error {
	task, err := m.getTask(taskId, userId)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrTaskNotStarted
	}
	if task != nil && task.FinishedAt != nil {
		return ErrTaskAlreadyFinished
	}
	query := `
		UPDATE tasks_users SET finished_at = $1 WHERE task_id = $2
	`

	args := []any{time.Now().UTC().Format(time.RFC3339), taskId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		slog.Error("error stopping task", "err", err.Error())
		return err
	}

	return nil
}

func (m TasksPostgres) getTask(taskId, userId string) (*Task, error) {
	query := `
		SELECT t.id, t.project_id, t.name, t.description, u.username, tu.started_at, tu.finished_at 
		FROM tasks t 
		INNER JOIN tasks_users tu ON t.id = tu.task_id
		INNER JOIN users u ON u.id = tu.user_id
		WHERE t.id = $1 AND u.id = $2;
	`

	args := []any{taskId, userId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	task := Task{}
	err := m.DB.QueryRowContext(ctx, query, args...).
		Scan(&task.Id, &task.ProjectId, &task.Name, &task.Description, &task.Username, &task.StartedAt, &task.FinishedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		default:
			return nil, err
		}
	}

	return &task, nil
}
