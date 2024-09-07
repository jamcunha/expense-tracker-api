package expense

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/database"
	"github.com/jamcunha/expense-tracker/internal/model"
	"github.com/shopspring/decimal"
)

type SqlcRepo struct {
	DB      *sql.DB
	Queries *database.Queries
}

func (s *SqlcRepo) Create(ctx context.Context, expense model.Expense) (model.Expense, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return model.Expense{}, err
	}
	defer tx.Rollback()

	qtx := s.Queries.WithTx(tx)

	dbExpense, err := qtx.CreateExpense(ctx, database.CreateExpenseParams{
		ID:          expense.ID,
		CreatedAt:   expense.CreatedAt,
		UpdatedAt:   expense.UpdatedAt,
		Description: expense.Description,
		Amount:      expense.Amount.String(),
		CategoryID:  expense.CategoryID,
		UserID:      expense.UserID,
	})
	if err != nil {
		return model.Expense{}, err
	}

	if err := qtx.UpdateBudgetAmount(ctx, database.UpdateBudgetAmountParams{
		CategoryID: dbExpense.CategoryID,
		Amount:     dbExpense.Amount,
		StartDate:  dbExpense.CreatedAt,
	}); err != nil {
		return model.Expense{}, err
	}

	return expense, tx.Commit()
}

func (s *SqlcRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Queries.WithTx(tx)

	dbExpense, err := qtx.DeleteExpense(ctx, id)
	if err != nil {
		return err
	}

	if err := qtx.UpdateBudgetAmount(ctx, database.UpdateBudgetAmountParams{
		CategoryID: dbExpense.CategoryID,
		Amount:     "-" + dbExpense.Amount,
		StartDate:  dbExpense.CreatedAt,
	}); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SqlcRepo) FindByID(ctx context.Context, id uuid.UUID) (model.Expense, error) {
	dbExpense, err := s.Queries.GetExpenseByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Expense{}, ErrNotFound
	} else if err != nil {
		return model.Expense{}, err
	}

	return model.Expense{
		ID:        dbExpense.ID,
		CreatedAt: dbExpense.CreatedAt,
		UpdatedAt: dbExpense.UpdatedAt,

		Description: dbExpense.Description,
		Amount:      decimal.RequireFromString(dbExpense.Amount),
		CategoryID:  dbExpense.CategoryID,
		UserID:      dbExpense.UserID,
	}, nil
}

func (s *SqlcRepo) FindByCategory(
	ctx context.Context,
	categoryID uuid.UUID,
	page FindAllPage,
) (FindResult, error) {
	var dbExpenses []database.Expense
	var err error

	if page.Cursor == "" {
		dbExpenses, err = s.Queries.GetCategoryExpenses(ctx, database.GetCategoryExpensesParams{
			CategoryID: categoryID,
			Limit:      page.Limit,
		})
	} else {
		t, id, err := decodeCursor(page.Cursor)
		if err != nil {
			return FindResult{}, err
		}

		dbExpenses, err = s.Queries.GetCategoryExpensesPaged(ctx, database.GetCategoryExpensesPagedParams{
			CategoryID: categoryID,
			CreatedAt:  t,
			ID:         id,
			Limit:      page.Limit,
		})
	}

	if errors.Is(err, sql.ErrNoRows) {
		return FindResult{}, ErrNotFound
	} else if err != nil {
		return FindResult{}, err
	}

	// NOTE: this may not be needed if empty arrays are not created and ErrNotFound is returned insted
	// (to test)
	if len(dbExpenses) == 0 {
		return FindResult{
			Expenses: []model.Expense{},
			Cursor:   "",
		}, nil
	}

	expenses := make([]model.Expense, len(dbExpenses))
	for i, dbExpense := range dbExpenses {
		expenses[i] = model.Expense{
			ID:        dbExpense.ID,
			CreatedAt: dbExpense.CreatedAt,
			UpdatedAt: dbExpense.UpdatedAt,

			Description: dbExpense.Description,
			Amount:      decimal.RequireFromString(dbExpense.Amount),
			CategoryID:  dbExpense.CategoryID,
			UserID:      dbExpense.UserID,
		}
	}

	cursor := ""
	if len(expenses) == int(page.Limit) {
		cursor = encodeCursor(expenses[len(expenses)-1].CreatedAt, expenses[len(expenses)-1].ID)
	}

	return FindResult{
		Expenses: expenses,
		Cursor:   cursor,
	}, nil
}

func (s *SqlcRepo) FindAll(
	ctx context.Context,
	userID uuid.UUID,
	page FindAllPage,
) (FindResult, error) {
	var dbExpenses []database.Expense
	var err error

	if page.Cursor == "" {
		dbExpenses, err = s.Queries.GetUserExpenses(ctx, database.GetUserExpensesParams{
			UserID: userID,
			Limit:  page.Limit,
		})
	} else {
		t, id, err := decodeCursor(page.Cursor)
		if err != nil {
			return FindResult{}, err
		}

		dbExpenses, err = s.Queries.GetUserExpensesPaged(ctx, database.GetUserExpensesPagedParams{
			UserID:    userID,
			CreatedAt: t,
			ID:        id,
			Limit:     page.Limit,
		})
	}

	if errors.Is(err, sql.ErrNoRows) {
		return FindResult{}, ErrNotFound
	} else if err != nil {
		return FindResult{}, err
	}

	// NOTE: this may not be needed if empty arrays are not created and ErrNotFound is returned insted
	// (to test)
	if len(dbExpenses) == 0 {
		return FindResult{
			Expenses: []model.Expense{},
			Cursor:   "",
		}, nil
	}

	expenses := make([]model.Expense, len(dbExpenses))
	for i, dbExpense := range dbExpenses {
		expenses[i] = model.Expense{
			ID:        dbExpense.ID,
			CreatedAt: dbExpense.CreatedAt,
			UpdatedAt: dbExpense.UpdatedAt,

			Description: dbExpense.Description,
			Amount:      decimal.RequireFromString(dbExpense.Amount),
			CategoryID:  dbExpense.CategoryID,
			UserID:      dbExpense.UserID,
		}
	}

	cursor := ""
	if len(expenses) == int(page.Limit) {
		cursor = encodeCursor(expenses[len(expenses)-1].CreatedAt, expenses[len(expenses)-1].ID)
	}

	return FindResult{
		Expenses: expenses,
		Cursor:   cursor,
	}, nil
}

func decodeCursor(encodedCursor string) (time.Time, uuid.UUID, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return time.Time{}, uuid.UUID{}, ErrInvalidCursor
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		return time.Time{}, uuid.UUID{}, ErrInvalidCursor
	}

	t, err := time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return time.Time{}, uuid.UUID{}, err
	}

	id, err := uuid.Parse(arrStr[1])
	if err != nil {
		return time.Time{}, uuid.UUID{}, err
	}

	return t, id, nil
}

func encodeCursor(t time.Time, uuid uuid.UUID) string {
	return base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), uuid.String()),
	))
}
