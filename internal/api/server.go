package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jamcunha/expense-tracker/internal/database"

	_ "github.com/lib/pq"
)

type Server struct {
	addr string
	DB   *database.Queries
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func NewServer(addr string, dbUrl string) *Server {
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}

	return &Server{
		addr: addr,
		DB:   database.New(conn),
	}
}

func (s *Server) Start() error {
	r := http.NewServeMux()

	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, struct{}{})
	})

	r.HandleFunc("POST /users", s.handleCreateUser)
	r.HandleFunc("DELETE /users", s.handleDeleteUser)

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", r))

	return http.ListenAndServe(s.addr, v1)
}

type ApiError struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
