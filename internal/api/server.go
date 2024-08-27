package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/jamcunha/expense-tracker/internal/database"
	"github.com/jamcunha/expense-tracker/internal/middleware"
	"github.com/jamcunha/expense-tracker/internal/utils"

	_ "github.com/lib/pq"
)

type Server struct {
	addr string
	DB   *database.Queries
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func NewServer() *Server {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}

	return &Server{
		addr: ":" + port,
		DB:   database.New(conn),
	}
}

func (s *Server) Start() error {
	r := http.NewServeMux()

	jwtMiddleware := func(f http.HandlerFunc) http.Handler { return middleware.JWTAuth(f) }

	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, struct{}{})
	})

	r.HandleFunc("POST /users", s.handleCreateUser)
	r.HandleFunc("DELETE /users", s.handleDeleteUser)
	r.HandleFunc("POST /login", s.handleUserLogin)

	r.Handle("GET /categories", jwtMiddleware(s.handleGetUserCategories))
	r.Handle("POST /categories", jwtMiddleware(s.handleCreateCategory))
	// maybe use DELETE /categories/:id
	r.Handle("DELETE /categories", jwtMiddleware(s.handleDeleteCategory))

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", r))

	log.Printf("server listening on %s", s.addr)
	return http.ListenAndServe(s.addr, v1)
}
