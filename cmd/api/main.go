package main

import (
	"flag"
	"log"

	"github.com/jamcunha/expense-tracker/internal/api"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	flag.Parse()

	server := api.NewServer(*addr)
	log.Printf("Server started on %s", *addr)
	log.Fatal(server.Start())
}
