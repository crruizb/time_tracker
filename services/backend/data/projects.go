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
