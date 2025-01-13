package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal"
	"github.com/jamcunha/expense-tracker/internal/repository"
	"github.com/jamcunha/expense-tracker/internal/service"
)

type Category struct {
	service service.Category
}

func NewCategory(db internal.DBConn, queries internal.Querier) *Category {
	return &Category{
		service: service.Category{
			DB:      db,
			Queries: queries,
		},
	}
}

func (h *Category) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	c, err := h.service.GetByID(r.Context(), id, userID)
	if errors.Is(err, service.ErrCategoryNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Category does not exist"}`))
		return
	} else if err != nil {
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
	// Default page limit
	limit := int32(10)

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		const decimal = 10
		const bitSize = 32
		limitParsed, err := strconv.ParseInt(limitStr, decimal, bitSize)
		if err != nil || limit < 1 {
			fmt.Println("Handler Error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		limit = int32(limitParsed)
	}

	cur := r.URL.Query().Get("cursor")

	userID := r.Context().Value("userID").(uuid.UUID)

	categories, err := h.service.GetAll(r.Context(), userID, limit, cur)
	if errors.Is(err, service.ErrCategoryNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "No categories found"}`))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Categories []repository.Category `json:"categories"`
		Next       string                `json:"next,omitempty"`
	}

	response.Categories = categories
	response.Next = ""

	if len(categories) == int(limit) {
		lastCategory := categories[len(categories)-1]
		cur = internal.EncodeCursor(lastCategory.CreatedAt, lastCategory.ID)
	}

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

	c, err := h.service.Create(r.Context(), body.Name, userID)
	if err != nil {
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

func (h *Category) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var body struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	userID := r.Context().Value("userID").(uuid.UUID)

	c, err := h.service.Update(r.Context(), id, userID, body.Name)
	if errors.Is(err, service.ErrCategoryNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Category does not exist"}`))
		return
	} else if err != nil {
		fmt.Println("failed to update:", err)
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
	w.WriteHeader(http.StatusOK)

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

	c, err := h.service.DeleteByID(r.Context(), id, userID)
	if errors.Is(err, service.ErrCategoryNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Category does not exist"}`))
		return
	} else if err != nil {
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
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}
