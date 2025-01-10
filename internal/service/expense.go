package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jamcunha/expense-tracker/internal/repository"
	"github.com/shopspring/decimal"
)

type Expense struct {
	DB      *pgx.Conn
	Queries *repository.Queries
}

func (s *Expense) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	e, err := s.Queries.GetExpenseByID(r.Context(), repository.GetExpenseByIDParams{
		ID:     id,
		UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
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

func (s *Expense) GetAll(w http.ResponseWriter, r *http.Request) {
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

	var expenses []repository.Expense

	if cursor == "" {
		expenses, err = s.Queries.GetUserExpenses(r.Context(), repository.GetUserExpensesParams{
			UserID: userID,
			Limit:  int32(limit),
		})
	} else {
		t, id, err := decodeCursor(cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		expenses, err = s.Queries.GetUserExpensesPaged(r.Context(), repository.GetUserExpensesPagedParams{
			UserID:    userID,
			CreatedAt: t,
			ID:        id,
			Limit:     int32(limit),
		})
	}

	if errors.Is(err, pgx.ErrNoRows) {
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
		Expenses []repository.Expense `json:"expenses"`
		Next     string               `json:"next,omitempty"`
	}

	response.Expenses = expenses
	response.Next = ""

	if len(expenses) == int(limit) {
		lastExpense := expenses[len(expenses)-1]
		response.Next = encodeCursor(lastExpense.CreatedAt, lastExpense.ID)
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

func (s *Expense) GetByCategory(w http.ResponseWriter, r *http.Request) {
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

	var expenses []repository.Expense

	if cursor == "" {
		expenses, err = s.Queries.GetCategoryExpenses(
			r.Context(),
			repository.GetCategoryExpensesParams{
				CategoryID: categoryID,
				UserID:     userID,
				Limit:      int32(limit),
			},
		)
	} else {
		t, id, err := decodeCursor(cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		expenses, err = s.Queries.GetCategoryExpensesPaged(r.Context(), repository.GetCategoryExpensesPagedParams{
			CategoryID: categoryID,
			UserID:     userID,
			CreatedAt:  t,
			ID:         id,
			Limit:      int32(limit),
		})
	}

	if errors.Is(err, pgx.ErrNoRows) {
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
		Expenses []repository.Expense `json:"expenses"`
		Next     string               `json:"next,omitempty"`
	}

	response.Expenses = expenses
	response.Next = ""

	if len(expenses) == int(limit) {
		lastExpenses := expenses[len(expenses)-1]
		response.Next = encodeCursor(lastExpenses.CreatedAt, lastExpenses.ID)
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

func (s *Expense) Create(w http.ResponseWriter, r *http.Request) {
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
	e, err := s.Queries.CreateExpense(r.Context(), repository.CreateExpenseParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,

		Description: body.Description,
		Amount:      decimal.NewFromFloat(body.Amount),
		CategoryID:  body.CategoryID,
		UserID:      userID,
	})
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

func (s *Expense) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		fmt.Println("Handler Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	tx, err := s.DB.Begin(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	qtx := s.Queries.WithTx(tx)

	e, err := qtx.DeleteExpense(r.Context(), repository.DeleteExpenseParams{
		ID:     id,
		UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Expense does not exist"}`))
		return
	} else if err != nil {
		fmt.Println("failed to delete:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = qtx.UpdateBudgetAmount(r.Context(), repository.UpdateBudgetAmountParams{
		CategoryID: e.CategoryID,
		Amount:     e.Amount.Neg(),
		StartDate:  e.CreatedAt,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(e)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if tx.Commit(r.Context()) != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}

func (s *Expense) Update(w http.ResponseWriter, r *http.Request) {
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

	tx, err := s.DB.Begin(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	qtx := s.Queries.WithTx(tx)

	e, err := qtx.GetExpenseByID(r.Context(), repository.GetExpenseByIDParams{
		ID:     id,
		UserID: userID,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`{"error": "Expense does not exist"}`))
		return
	} else if err != nil {
		fmt.Println("failed to update:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	oldCategory := e.CategoryID
	oldAmount := e.Amount

	if body.Description == "" {
		body.Description = e.Description
	}

	if bodyAmount.IsZero() {
		bodyAmount = e.Amount
	}

	if categoryID == uuid.Nil {
		categoryID = e.CategoryID
	}

	now := time.Now()
	e, err = qtx.UpdateExpense(r.Context(), repository.UpdateExpenseParams{
		ID:          id,
		UserID:      userID,
		Description: body.Description,
		Amount:      bodyAmount,
		CategoryID:  categoryID,
		UpdatedAt:   now,
	})
	if err != nil {
		fmt.Println("failed to update:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = qtx.UpdateBudgetAmount(r.Context(), repository.UpdateBudgetAmountParams{
		CategoryID: oldCategory,
		Amount:     oldAmount.Neg(),
		StartDate:  e.CreatedAt,
	})
	if err != nil {
		fmt.Println("failed to update:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = qtx.UpdateBudgetAmount(r.Context(), repository.UpdateBudgetAmountParams{
		CategoryID: e.CategoryID,
		Amount:     e.Amount,
		StartDate:  e.UpdatedAt, // Should this be CreatedAt?
	})
	if err != nil {
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

	if tx.Commit(r.Context()) != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}
