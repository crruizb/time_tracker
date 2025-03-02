package data

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
)

type UsersStore struct {
	DB *sql.DB
}

func NewUsersStore(db *sql.DB) *UsersStore {
	return &UsersStore{
		DB: db,
	}
}

type User struct {
	Id       string `db:"id"`
	Username string `json:"login" db:"username"`
	Source   string `db:"source"`
}

func (m UsersStore) GetUser(username, source string) (*User, error) {
	query := `
		SELECT id, username, source FROM users
		WHERE username = $1 AND source = $2
	`

	args := []any{username, source}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &User{}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Username, &user.Source)
	if err != nil {
		slog.Error("error finding user", "err", err.Error())
		return nil, err
	}

	return user, nil
}

func (m UsersStore) InsertUser(username, source string) (*User, error) {
	query := `
		INSERT INTO users (username, source)
		VALUES ($1, $2)
		RETURNING id
	`

	args := []any{username, source}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &User{
		Username: username,
		Source:   source,
	}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id)
	if err != nil {
		slog.Error("error creating user", "err", err.Error())
		return nil, err
	}

	return user, nil
}
