package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/jamcunha/expense-tracker/internal/application"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed to load env:", err)
		return
	}

	cfg, err := application.LoadConfig()
	if err != nil {
		fmt.Println("failed to load config:", err)
		return
	}

	app, err := application.New(cfg)
	if err != nil {
		fmt.Println("failed to create application:", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start application:", err)
	}
}
