package main

import (
	"log"

	"github.com/jamcunha/expense-tracker/internal/api"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	server := api.NewServer()
	log.Fatal(server.Start())
}
