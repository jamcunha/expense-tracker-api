package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/database"
	"github.com/jamcunha/expense-tracker/internal/model"
	"github.com/jamcunha/expense-tracker/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			utils.ApiError{Error: fmt.Sprintf("Invalid JSON: %v", err)},
		)
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(params.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			utils.ApiError{Error: fmt.Sprintf("Error encrypting password: %v", err)},
		)
	}

	user, err := s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Email:     params.Email,
		Password:  string(encryptedPassword),
	})
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			utils.ApiError{Error: fmt.Sprintf("Error creating user: %v", err)},
		)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, model.DatabaseUserToUser(user))
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID string `json:"id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			utils.ApiError{Error: fmt.Sprintf("Invalid JSON: %v", err)},
		)
		return
	}

	id, err := uuid.Parse(params.ID)
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			utils.ApiError{Error: fmt.Sprintf("Invalid ID: %v", err)},
		)
		return
	}

	err = s.DB.DeleteUser(r.Context(), id)
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			utils.ApiError{Error: fmt.Sprintf("Error deleting user: %v", err)},
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			utils.ApiError{Error: fmt.Sprintf("Invalid JSON: %v", err)},
		)
		return
	}

	dbUser, err := s.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(
				w,
				http.StatusUnauthorized,
				utils.ApiError{Error: "Invalid credentials"},
			)
			return
		}

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			utils.ApiError{Error: fmt.Sprintf("Error logging in: %v", err)},
		)
		return
	}

	user := model.DatabaseUserToUser(dbUser)

	if !user.ComparePassword(params.Password) {
		utils.WriteJSON(
			w,
			http.StatusUnauthorized,
			utils.ApiError{Error: "Invalid credentials"},
		)
		return
	}

	token, err := createJWT(user)
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			utils.ApiError{Error: fmt.Sprintf("Error creating token: %v", err)},
		)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func createJWT(user model.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "expense-tracker",
		Subject:   user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
