package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jamcunha/expense-tracker/internal"
	"github.com/jamcunha/expense-tracker/internal/middleware"
	"github.com/jamcunha/expense-tracker/internal/repository"

	"github.com/jackc/pgx/v5"
)

type App struct {
	router *http.ServeMux

	DB      internal.DBConn
	Queries internal.Querier
	config  Config
}

func New(config Config) (*App, error) {
	conn, err := pgx.Connect(context.Background(), config.PostgresUrl)
	if err != nil {
		return &App{}, fmt.Errorf("error opening database connection: %w", err)
	}

	app := &App{
		DB:      conn,
		Queries: internal.NewQuerier(repository.New(conn)),
		config:  config,
	}
	app.loadRoutes("/api/v1")

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

		a.DB.Close(ctx)
		return server.Shutdown(timeout)
	case err := <-ch:
		a.DB.Close(ctx)
		return err
	}
}
