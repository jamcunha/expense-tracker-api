package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/database"
	"github.com/jamcunha/expense-tracker/internal/model"
)

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// TODO: use bcrypt2

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ApiError{Error: fmt.Sprintf("Invalid JSON: %v", err)})
	}

	user, err := s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Email:     params.Email,
		Password:  params.Password,
	})
	if err != nil {
		WriteJSON(
			w,
			http.StatusInternalServerError,
			ApiError{Error: fmt.Sprintf("Error creating user: %v", err)},
		)
	}

	WriteJSON(w, http.StatusCreated, model.DatabaseUserToUser(user))
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID string `json:"id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ApiError{Error: fmt.Sprintf("Invalid JSON: %v", err)})
	}

	id, err := uuid.Parse(params.ID)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ApiError{Error: fmt.Sprintf("Invalid ID: %v", err)})
	}

	err = s.DB.DeleteUser(r.Context(), id)
	if err != nil {
		WriteJSON(
			w,
			http.StatusInternalServerError,
			ApiError{Error: fmt.Sprintf("Error deleting user: %v", err)},
		)
	}

	w.WriteHeader(http.StatusNoContent)
}
