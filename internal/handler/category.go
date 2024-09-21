package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/model"
	repo "github.com/jamcunha/expense-tracker/internal/repository/category"
)

type Category struct {
	Repo repo.Repo
}

func (h *Category) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	c, err := h.Repo.FindByID(r.Context(), id, userID)
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Category does not exist"}`))
		return
	} else if err != nil {
		fmt.Print("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(c)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}

func (h *Category) GetAll(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	const decimal = 10
	const bitSize = 32
	limit, err := strconv.ParseInt(limitStr, decimal, bitSize)
	if err != nil || limit < 1 {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cursor := r.URL.Query().Get("cursor")

	userID := r.Context().Value("userID").(uuid.UUID)

	categories, err := h.Repo.FindAll(r.Context(), userID, repo.FindAllPage{
		Limit:  int32(limit),
		Cursor: cursor,
	})
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "No categories found"}`))
		return
	} else if err != nil {
		fmt.Println("failed to find:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Categories []model.Category `json:"categories"`
		Next       string           `json:"next,omitempty"`
	}

	response.Categories = categories.Categories
	response.Next = categories.Cursor

	res, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}

func (h *Category) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	now := time.Now()
	c := model.Category{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      body.Name,
		UserID:    userID,
	}

	c, err := h.Repo.Create(r.Context(), c)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(c)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	w.Write(res)
}

func (h *Category) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Hander Error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	user, err := h.Repo.Delete(r.Context(), id, userID)
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Category does not exist"}`))
		return
	} else if err != nil {
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
