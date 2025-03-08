package data

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
)

type ProjectsPostgres struct {
	DB *sql.DB
}

func NewProjectsPostgres(db *sql.DB) *ProjectsPostgres {
	return &ProjectsPostgres{
		DB: db,
	}
}

type Project struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectTask struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Tasks       []TaskOfProject `json:"tasks"`
}

type TaskOfProject struct {
	Id          *string    `json:"id" db:"id"`
	Name        *string    `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	Username    *string    `json:"username" db:"username"`
	StartedAt   *time.Time `json:"startedAt" db:"started_at"`
	FinishedAt  *time.Time `json:"finishedAt" db:"finished_at"`
}

func (m ProjectsPostgres) CreateProject(name, description, userId string) (*Project, error) {
	query := `
		INSERT INTO projects (name, description, user_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	args := []any{name, description, userId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	project := &Project{
		Name:        name,
		Description: description,
	}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&project.Id)
	if err != nil {
		slog.Error("error creating project", "err", err.Error())
		return nil, err
	}

	return project, nil
}

func (m ProjectsPostgres) GetProjects(userId string) ([]ProjectTask, error) {
	query := `
		SELECT p.id, p.name, p.description, t.id as taskId, t.name as taskName, t.description as taskDescription, u.username, tu.started_at, tu.finished_at
		FROM projects p
		LEFT JOIN tasks t ON p.id = t.project_id 
		LEFT JOIN tasks_users tu ON t.id = tu.task_id
		LEFT JOIN users u ON u.id = tu.user_id
		WHERE p.user_id = $1; 
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	projects := map[string]ProjectTask{}
	for rows.Next() {
		var (
			id              string
			name            string
			description     string
			taskId          *string
			taskName        *string
			taskDescription *string
			username        *string
			started_at      *time.Time
			finished_at     *time.Time
		)

		err := rows.Scan(&id, &name, &description, &taskId, &taskName, &taskDescription, &username, &started_at, &finished_at)
		if err != nil {
			return nil, err
		}

		savedProjects, ok := projects[id]
		task := TaskOfProject{
			Id:          taskId,
			Name:        taskName,
			Description: taskDescription,
			Username:    username,
			StartedAt:   started_at,
			FinishedAt:  finished_at,
		}

		if !ok {
			projectTask := ProjectTask{
				Id:          id,
				Name:        name,
				Description: description,
				Tasks:       []TaskOfProject{},
			}
			if task.Id != nil {
				projectTask.Tasks = append(projectTask.Tasks, task)
			}
			projects[id] = projectTask

		} else {
			savedProjects.Tasks = append(savedProjects.Tasks, task)
			projects[id] = savedProjects
		}
	}

	projectTasks := []ProjectTask{}
	for _, v := range projects {
		projectTasks = append(projectTasks, v)
	}

	return projectTasks, nil
}
