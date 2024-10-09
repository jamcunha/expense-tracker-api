package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/model"
	repo "github.com/jamcunha/expense-tracker/internal/repository/user"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Repo repo.Repo
}

func (h *User) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.Repo.FindByID(r.Context(), id)
	if errors.Is(err, repo.ErrNotFound) {
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

	u, err = h.Repo.Create(r.Context(), u)
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

	user, err := h.Repo.Delete(r.Context(), id)
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "User does not exist"}`))
		return
	}
	if err != nil {
		fmt.Println("failed to delete:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}
