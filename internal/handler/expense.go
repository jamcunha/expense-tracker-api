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
	"github.com/shopspring/decimal"
)

type Expense struct {
	service service.Expense
}

func NewExpense(db internal.DBConn, queries internal.Querier) *Expense {
	return &Expense{
		service: service.Expense{
			DB:      db,
			Queries: queries,
		},
	}
}

func (h *Expense) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	e, err := h.service.GetByID(r.Context(), id, userID)
	if errors.Is(err, service.ErrExpenseNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Expense does not exist"}`))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(e)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}

func (h *Expense) GetAll(w http.ResponseWriter, r *http.Request) {
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

	expenses, err := h.service.GetAll(r.Context(), userID, limit, cur)
	if errors.Is(err, service.ErrExpenseNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "No expenses found"}`))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Expenses []repository.Expense `json:"expenses"`
		Next     string               `json:"next,omitempty"`
	}

	response.Expenses = expenses
	response.Next = ""

	if len(expenses) == int(limit) {
		lastExpense := expenses[len(expenses)-1]
		response.Next = internal.EncodeCursor(lastExpense.CreatedAt, lastExpense.ID)
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

func (h *Expense) GetByCategory(w http.ResponseWriter, r *http.Request) {
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

	categoryID, err := uuid.Parse(r.PathValue("category_id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	expenses, err := h.service.GetByCategory(r.Context(), categoryID, userID, limit, cur)
	if errors.Is(err, service.ErrExpenseNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "No expenses found"}`))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Expenses []repository.Expense `json:"expenses"`
		Next     string               `json:"next,omitempty"`
	}

	response.Expenses = expenses
	response.Next = ""

	if len(expenses) == int(limit) {
		lastExpenses := expenses[len(expenses)-1]
		response.Next = internal.EncodeCursor(lastExpenses.CreatedAt, lastExpenses.ID)
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

func (h *Expense) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		CategoryID  string  `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	categoryID, err := uuid.Parse(r.PathValue("category_id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	e, err := h.service.Create(
		r.Context(),
		userID,
		body.Description,
		decimal.NewFromFloat(body.Amount),
		categoryID,
	)

	res, err := json.Marshal(e)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	w.Write(res)
}

func (h *Expense) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	e, err := h.service.DeleteByID(r.Context(), id, userID)
	if errors.Is(err, service.ErrExpenseNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Expense does not exist"}`))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(e)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}

func (h *Expense) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var body struct {
		Description string  `json:"description,omitempty"`
		Amount      float64 `json:"amount,omitempty"`
		CategoryID  string  `json:"category_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyAmount := decimal.NewFromFloat(body.Amount)

	categoryID, err := uuid.Parse(body.CategoryID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		w.Write([]byte(`{"error": "Invalid category ID"}`))
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	e, err := h.service.Update(r.Context(), id, categoryID, userID, body.Description, bodyAmount)
	if errors.Is(err, service.ErrExpenseNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Expense does not exist"}`))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(e)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}
