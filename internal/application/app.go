package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jamcunha/expense-tracker/internal/database"
	"github.com/jamcunha/expense-tracker/internal/middleware"

	_ "github.com/lib/pq" // find if there is a way to add this by default to sqlc generated files
)

type App struct {
	router *http.ServeMux

	DB      *sql.DB
	Queries *database.Queries
	config  Config
}

func New(config Config) (*App, error) {
	conn, err := sql.Open("postgres", *config.PostgresUrl)
	if err != nil {
		return &App{}, fmt.Errorf("error opening database connection: %w", err)
	}

	app := &App{
		DB:      conn,
		Queries: database.New(conn),
		config:  config,
	}
	app.loadV1Routes("/api/v1")

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	server := http.Server{
		Addr:    ":" + a.config.ServerPort,
		Handler: middleware.Logging(a.router),
	}

	ch := make(chan error, 1)

	go func() {
		fmt.Println("Server is running on port", a.config.ServerPort)
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	case err := <-ch:
		return err
	}
}
