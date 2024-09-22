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
	repo "github.com/jamcunha/expense-tracker/internal/repository/expense"
	"github.com/shopspring/decimal"
)

type Expense struct {
	Repo repo.Repo
}

func (h *Expense) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	e, err := h.Repo.FindByID(r.Context(), id, userID)
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Expense does not exist"}`))
		return
	} else if err != nil {
		fmt.Print("failed to insert:", err)
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

	expenses, err := h.Repo.FindAll(r.Context(), userID, repo.FindAllPage{
		Limit:  int32(limit),
		Cursor: cursor,
	})
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "No expenses found"}`))
		return
	} else if err != nil {
		fmt.Println("failed to find:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Expenses []model.Expense `json:"expenses"`
		Next     string          `json:"next,omitempty"`
	}

	response.Expenses = expenses.Expenses
	response.Next = expenses.Cursor

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

	categoryID, err := uuid.Parse(r.PathValue("category_id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	expenses, err := h.Repo.FindByCategory(r.Context(), categoryID, userID, repo.FindAllPage{
		Limit:  int32(limit),
		Cursor: cursor,
	})
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "No expenses found"}`))
		return
	} else if err != nil {
		fmt.Println("failed to find:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Expenses []model.Expense `json:"expenses"`
		Next     string          `json:"next,omitempty"`
	}

	response.Expenses = expenses.Expenses
	response.Next = expenses.Cursor

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
		Description string    `json:"description"`
		Amount      float64   `json:"amount"`
		CategoryID  uuid.UUID `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	now := time.Now()
	e := model.Expense{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,

		Description: body.Description,
		Amount:      decimal.NewFromFloat(body.Amount),
		CategoryID:  body.CategoryID,
		UserID:      userID,
	}

	e, err := h.Repo.Create(r.Context(), e)
	if err != nil {
		fmt.Println("failed to insert:", err)
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

	expense, err := h.Repo.Delete(r.Context(), id, userID)
	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Expense does not exist"}`))
		return
	} else if err != nil {
		fmt.Println("failed to delete:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(expense)
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

	categoryID, err := uuid.Parse(body.CategoryID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		w.Write([]byte(`{"error": "Invalid category ID"}`))
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	e, err := h.Repo.Update(r.Context(), repo.UpdateExpense{
		ID:     id,
		UserID: userID,

		Description: body.Description,
		Amount:      decimal.NewFromFloat(body.Amount),
		CategoryID:  categoryID,
	})

	if errors.Is(err, repo.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Expense does not exist"}`))
		return
	} else if err != nil {
		fmt.Println("failed to update:", err)
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
