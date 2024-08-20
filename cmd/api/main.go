package main

import (
	"flag"
	"log"

	"github.com/jamcunha/expense-tracker/internal/api"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	dbUrl := flag.String(
		"db-url",
		"postgres://postgres:postgres@localhost:5432/local-db?sslmode=disable",
		"database url",
	)

	flag.Parse()

	server := api.NewServer(*addr, *dbUrl)
	log.Printf("Server started on %s", *addr)
	log.Fatal(server.Start())
}
