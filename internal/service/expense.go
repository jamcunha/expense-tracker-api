package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jamcunha/expense-tracker/internal"
	"github.com/jamcunha/expense-tracker/internal/repository"
	"github.com/shopspring/decimal"
)

type Expense struct {
	DB      *pgx.Conn
	Queries *repository.Queries
}

func (s *Expense) GetByID(ctx context.Context, id, userID uuid.UUID) (repository.Expense, error) {
	e, err := s.Queries.GetExpenseByID(ctx, repository.GetExpenseByIDParams{
		ID:     id,
		UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.Expense{}, ErrExpenseNotFound
	} else if err != nil {
		fmt.Print("failed to insert:", err)
		return repository.Expense{}, err
	}

	return e, nil
}

func (s *Expense) GetAll(
	ctx context.Context,
	userID uuid.UUID,
	limit int32,
	cur string,
) ([]repository.Expense, error) {
	var expenses []repository.Expense
	var err error

	if cur == "" {
		expenses, err = s.Queries.GetUserExpenses(ctx, repository.GetUserExpensesParams{
			UserID: userID,
			Limit:  int32(limit),
		})
	} else {
		t, id, err := internal.DecodeCursor(cur)
		if err != nil {
			return []repository.Expense{}, ErrDecodeCursor
		}

		expenses, err = s.Queries.GetUserExpensesPaged(ctx, repository.GetUserExpensesPagedParams{
			UserID:    userID,
			CreatedAt: t,
			ID:        id,
			Limit:     int32(limit),
		})
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return []repository.Expense{}, ErrExpenseNotFound
	} else if err != nil {
		fmt.Println("failed to find:", err)
		return []repository.Expense{}, err
	}

	return expenses, nil
}

func (s *Expense) GetByCategory(
	ctx context.Context,
	categoryID, userID uuid.UUID,
	limit int32,
	cur string,
) ([]repository.Expense, error) {
	var expenses []repository.Expense
	var err error

	if cur == "" {
		expenses, err = s.Queries.GetCategoryExpenses(
			ctx,
			repository.GetCategoryExpensesParams{
				CategoryID: categoryID,
				UserID:     userID,
				Limit:      int32(limit),
			},
		)
	} else {
		t, id, err := internal.DecodeCursor(cur)
		if err != nil {
			return []repository.Expense{}, ErrDecodeCursor
		}

		expenses, err = s.Queries.GetCategoryExpensesPaged(ctx, repository.GetCategoryExpensesPagedParams{
			CategoryID: categoryID,
			UserID:     userID,
			CreatedAt:  t,
			ID:         id,
			Limit:      int32(limit),
		})
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return []repository.Expense{}, ErrExpenseNotFound
	} else if err != nil {
		fmt.Println("failed to find:", err)
		return []repository.Expense{}, err
	}

	return expenses, nil
}

func (s *Expense) Create(
	ctx context.Context,
	userID uuid.UUID,
	description string,
	amount decimal.Decimal,
	categoryID uuid.UUID,
) (repository.Expense, error) {
	now := time.Now()
	e, err := s.Queries.CreateExpense(ctx, repository.CreateExpenseParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,

		Description: description,
		Amount:      amount,
		CategoryID:  categoryID,
		UserID:      userID,
	})
	if err != nil {
		fmt.Println("failed to insert:", err)
		return repository.Expense{}, err
	}

	return e, nil
}

func (s *Expense) DeleteByID(
	ctx context.Context,
	id, userID uuid.UUID,
) (repository.Expense, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return repository.Expense{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.Queries.WithTx(tx)

	e, err := qtx.DeleteExpense(ctx, repository.DeleteExpenseParams{
		ID:     id,
		UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.Expense{}, ErrExpenseNotFound
	} else if err != nil {
		fmt.Println("failed to delete:", err)
		return repository.Expense{}, err
	}

	err = qtx.UpdateBudgetAmount(ctx, repository.UpdateBudgetAmountParams{
		CategoryID: e.CategoryID,
		Amount:     e.Amount.Neg(),
		StartDate:  e.CreatedAt,
	})
	if err != nil {
		return repository.Expense{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.Expense{}, err
	}

	return e, nil
}

func (s *Expense) Update(
	ctx context.Context,
	id, categoryID, userID uuid.UUID,
	description string,
	amount decimal.Decimal,
) (repository.Expense, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return repository.Expense{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.Queries.WithTx(tx)

	e, err := qtx.GetExpenseByID(ctx, repository.GetExpenseByIDParams{
		ID:     id,
		UserID: userID,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return repository.Expense{}, ErrExpenseNotFound
	} else if err != nil {
		fmt.Println("failed to update:", err)
		return repository.Expense{}, err
	}

	oldCategory := e.CategoryID
	oldAmount := e.Amount

	if description == "" {
		description = e.Description
	}

	if amount.IsZero() {
		amount = e.Amount
	}

	if categoryID == uuid.Nil {
		categoryID = e.CategoryID
	}

	now := time.Now()
	e, err = qtx.UpdateExpense(ctx, repository.UpdateExpenseParams{
		ID:          id,
		UserID:      userID,
		Description: description,
		Amount:      amount,
		CategoryID:  categoryID,
		UpdatedAt:   now,
	})
	if err != nil {
		fmt.Println("failed to update:", err)
		return repository.Expense{}, err
	}

	err = qtx.UpdateBudgetAmount(ctx, repository.UpdateBudgetAmountParams{
		CategoryID: oldCategory,
		Amount:     oldAmount.Neg(),
		StartDate:  e.CreatedAt,
	})
	if err != nil {
		fmt.Println("failed to update:", err)
		return repository.Expense{}, err
	}

	err = qtx.UpdateBudgetAmount(ctx, repository.UpdateBudgetAmountParams{
		CategoryID: e.CategoryID,
		Amount:     e.Amount,
		StartDate:  e.UpdatedAt, // Should this be CreatedAt?
	})
	if err != nil {
		fmt.Println("failed to update:", err)
		return repository.Expense{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.Expense{}, err
	}

	return e, nil
}
