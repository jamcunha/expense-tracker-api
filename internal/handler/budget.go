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
	repo "github.com/jamcunha/expense-tracker/internal/repository/budget"
	"github.com/shopspring/decimal"
)

type Budget struct {
	Repo repo.Repo
}

func (h *Budget) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	b, err := h.Repo.FindByID(r.Context(), id, userID)
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Budget does not exist"}`))
		return
	} else if err != nil {
		fmt.Print("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(b)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}

func (h *Budget) GetAll(w http.ResponseWriter, r *http.Request) {
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

	budgets, err := h.Repo.FindAll(r.Context(), userID, repo.FindAllPage{
		Limit:  int32(limit),
		Cursor: cursor,
	})
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "No budgets found"}`))
		return
	} else if err != nil {
		fmt.Print("failed to find:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Budgets []model.Budget `json:"budgets"`
		Next    string         `json:"next,omitempty"`
	}

	response.Budgets = budgets.Budgets
	response.Next = budgets.Cursor

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

func (h *Budget) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Goal       float64 `json:"goal"`
		StartDate  string  `json:"start_date"`
		EndDate    string  `json:"end_date"`
		CategoryID string  `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)
	categoryID, err := uuid.Parse(body.CategoryID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		w.Write([]byte(`{"error": "Invalid category ID"}`))
		return
	}

	startDate, err := time.Parse(time.DateOnly, body.StartDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		w.Write([]byte(`{"error": "Invalid date format. Use YYYY-MM-DD"}`))
		return
	}

	endDate, err := time.Parse(time.DateOnly, body.EndDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		w.Write([]byte(`{"error": "Invalid date format. Use YYYY-MM-DD"}`))
		return
	}

	now := time.Now()
	b := model.Budget{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,

		Amount:     decimal.NewFromInt(0),
		Goal:       decimal.NewFromFloat(body.Goal),
		StartDate:  startDate,
		EndDate:    endDate,
		UserID:     userID,
		CategoryID: categoryID,
	}

	b, err = h.Repo.Create(r.Context(), b)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(b)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	w.Write(res)
}

func (h *Budget) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	budget, err := h.Repo.Delete(r.Context(), id, userID)
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Budget does not exist"}`))
		return
	} else if err != nil {
		fmt.Println("failed to delete:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(budget)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}
