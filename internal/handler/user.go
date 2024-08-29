package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/model"
	"github.com/jamcunha/expense-tracker/internal/repository/user"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Repo user.Repo
}

func (h *User) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.Repo.FindByID(r.Context(), id)
	if errors.Is(err, user.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "User does not exist"}`))
		return
	} else if err != nil {
		fmt.Print("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(u)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}

func (h *User) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(body.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		fmt.Println("failed to encrypt password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	now := time.Now()
	u := model.User{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      body.Name,
		Email:     body.Email,
		Password:  string(encryptedPassword),
	}

	err = h.Repo.Create(r.Context(), u)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(u)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	w.Write(res)
}

func (h *User) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Repo.Delete(r.Context(), id)
	if err != nil {
		fmt.Println("failed to delete:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *User) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.Repo.FindByEmail(r.Context(), body.Email)
	if errors.Is(err, user.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte(`{"error": "Invalid credentials"}`))
		return
	} else if err != nil {
		fmt.Print("failed to query:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !u.ComparePassword(body.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte(`{"error": "Invalid credentials"}`))
		return
	}

	token, err := createJWT(u)
	if err != nil {
		fmt.Println("failed to create JWT:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
}

func createJWT(user model.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	// TODO: create a config and load all environment variables before starting the server
	if jwtSecret == "" {
		panic("JWT_SECRET is not set")
	}

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
