package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal"
	"github.com/jamcunha/expense-tracker/internal/repository"
	"github.com/jamcunha/expense-tracker/internal/service"
	"github.com/shopspring/decimal"
)

type Budget struct {
	service service.Budget
}

func NewBudget(db internal.DBConn, queries internal.Querier) *Budget {
	return &Budget{
		service: service.Budget{
			DB:      db,
			Queries: queries,
		},
	}
}

func (h *Budget) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	b, err := h.service.GetByID(r.Context(), id, userID)
	if errors.Is(err, service.ErrBudgetNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Budget does not exist"}`))
		return
	} else if err != nil {
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

	budgets, err := h.service.GetAll(r.Context(), userID, limit, cur)
	if errors.Is(err, service.ErrBudgetNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "No budgets found"}`))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Budgets []repository.Budget `json:"budgets"`
		Next    string              `json:"next,omitempty"`
	}

	response.Budgets = budgets
	response.Next = ""

	if len(budgets) == int(limit) {
		lastBudget := budgets[len(budgets)-1]
		response.Next = internal.EncodeCursor(lastBudget.CreatedAt, lastBudget.ID)
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

	b, err := h.service.Create(
		r.Context(),
		userID,
		categoryID,
		decimal.NewFromFloat(body.Goal),
		startDate,
		endDate,
	)
	if errors.Is(err, service.ErrCategoryNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Category does not exist"}`))
		return
	} else if err != nil {
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

	b, err := h.service.DeleteByID(r.Context(), id, userID)
	if errors.Is(err, service.ErrBudgetNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Budget does not exist"}`))
		return
	} else if err != nil {
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
