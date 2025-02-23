package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/crruizb/api"
	"github.com/crruizb/data"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgresql://user:user@localhost:5432/time_tracker?sslmode=disable"
	db, err := openDB(dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	models := data.NewProjectsPostgres(db)

	srv := api.NewServer(":8080", models)
	if err := srv.Run(); err != nil {
		panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}