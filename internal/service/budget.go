package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	cursor "github.com/jamcunha/expense-tracker/internal"
	"github.com/jamcunha/expense-tracker/internal/repository"
	"github.com/shopspring/decimal"
)

type Budget struct {
	DB      *pgx.Conn
	Queries *repository.Queries
}

func (s *Budget) GetByID(ctx context.Context, id, userID uuid.UUID) (repository.Budget, error) {
	b, err := s.Queries.GetBudgetByID(ctx, repository.GetBudgetByIDParams{
		ID:     id,
		UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.Budget{}, ErrBudgetNotFound
	} else if err != nil {
		fmt.Print("failed to insert:", err)
		return repository.Budget{}, err
	}

	return b, nil
}

func (s *Budget) GetAll(
	ctx context.Context,
	userID uuid.UUID,
	limit int32,
	cur string,
) ([]repository.Budget, error) {
	var budgets []repository.Budget
	var err error

	if cur == "" {
		budgets, err = s.Queries.GetUserBudgets(ctx, repository.GetUserBudgetsParams{
			UserID: userID,
			Limit:  int32(limit),
		})
	} else {
		t, id, err := cursor.DecodeCursor(cur)
		if err != nil {
			return []repository.Budget{}, ErrDecodeCursor
		}

		budgets, err = s.Queries.GetUserBudgetsPaged(ctx, repository.GetUserBudgetsPagedParams{
			UserID:    userID,
			CreatedAt: t,
			ID:        id,
			Limit:     int32(limit),
		})
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return []repository.Budget{}, ErrBudgetNotFound
	} else if err != nil {
		fmt.Print("failed to find:", err)
		return []repository.Budget{}, err
	}

	return budgets, nil
}

func (s *Budget) Create(
	ctx context.Context,
	userID, categoryID uuid.UUID,
	goal decimal.Decimal,
	startDate, endDate time.Time,
) (repository.Budget, error) {
	now := time.Now()
	budgetParams := repository.CreateBudgetParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,

		Amount:     decimal.Zero,
		Goal:       goal,
		StartDate:  startDate,
		EndDate:    endDate,
		UserID:     userID,
		CategoryID: categoryID,
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return repository.Budget{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.Queries.WithTx(tx)

	amount, err := qtx.GetTotalSpentInCategory(
		ctx,
		repository.GetTotalSpentInCategoryParams{
			UserID:     userID,
			CategoryID: categoryID,
			StartDate:  startDate,
			EndDate:    endDate,
		},
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.Budget{}, ErrCategoryNotFound
	} else if err != nil {
		return repository.Budget{}, err
	}

	budgetParams.Amount = amount

	b, err := qtx.CreateBudget(ctx, budgetParams)
	if err != nil {
		fmt.Println("failed to insert:", err)
		return repository.Budget{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.Budget{}, err
	}

	return b, nil
}

func (s *Budget) DeleteByID(ctx context.Context, id, userID uuid.UUID) (repository.Budget, error) {
	b, err := s.Queries.DeleteBudget(ctx, repository.DeleteBudgetParams{
		ID:     id,
		UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.Budget{}, ErrBudgetNotFound
	} else if err != nil {
		fmt.Println("failed to delete:", err)
		return repository.Budget{}, err
	}

	return b, nil
}
