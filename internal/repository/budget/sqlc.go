package budget

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

func (s *SqlcRepo) Create(ctx context.Context, budget model.Budget) (model.Budget, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return model.Budget{}, err
	}
	defer tx.Rollback()

	qtx := s.Queries.WithTx(tx)

	amount, err := qtx.GetTotalSpentInCategory(ctx, database.GetTotalSpentInCategoryParams{
		UserID:     budget.UserID,
		CategoryID: budget.CategoryID,

		CreatedAt:   budget.StartDate,
		CreatedAt_2: budget.EndDate,
	})
	if err != nil {
		return model.Budget{}, err
	}

	_, err = qtx.CreateBudget(ctx, database.CreateBudgetParams{
		ID:        budget.ID,
		CreatedAt: budget.CreatedAt,
		UpdatedAt: budget.UpdatedAt,

		Amount:     amount,
		Goal:       budget.Goal.String(),
		StartDate:  budget.StartDate,
		EndDate:    budget.EndDate,
		UserID:     budget.UserID,
		CategoryID: budget.CategoryID,
	})
	if err != nil {
		return model.Budget{}, err
	}

	budget.Amount, err = decimal.NewFromString(amount)
	if err != nil {
		return model.Budget{}, err
	}

	return budget, tx.Commit()
}

func (s *SqlcRepo) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := s.Queries.DeleteBudget(ctx, database.DeleteBudgetParams{
		ID:     id,
		UserID: userID,
	})

	return err
}

func (s *SqlcRepo) FindByID(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) (model.Budget, error) {
	dbBudget, err := s.Queries.GetBudgetByID(ctx, database.GetBudgetByIDParams{
		ID:     id,
		UserID: userID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return model.Budget{}, ErrNotFound
	} else if err != nil {
		return model.Budget{}, err
	}

	return model.Budget{
		ID:        dbBudget.ID,
		CreatedAt: dbBudget.CreatedAt,
		UpdatedAt: dbBudget.UpdatedAt,

		Amount:     decimal.RequireFromString(dbBudget.Amount),
		Goal:       decimal.RequireFromString(dbBudget.Goal),
		StartDate:  dbBudget.StartDate,
		EndDate:    dbBudget.EndDate,
		UserID:     dbBudget.UserID,
		CategoryID: dbBudget.CategoryID,
	}, nil
}

func (s *SqlcRepo) FindAll(
	ctx context.Context,
	userID uuid.UUID,
	page FindAllPage,
) (FindResult, error) {
	var dbBudgets []database.Budget
	var err error

	if page.Cursor == "" {
		dbBudgets, err = s.Queries.GetUserBudgets(ctx, database.GetUserBudgetsParams{
			UserID: userID,
			Limit:  page.Limit,
		})
	} else {
		t, id, err := decodeCursor(page.Cursor)
		if err != nil {
			return FindResult{}, err
		}

		dbBudgets, err = s.Queries.GetUserBudgetsPaged(ctx, database.GetUserBudgetsPagedParams{
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
	if len(dbBudgets) == 0 {
		return FindResult{
			Budgets: []model.Budget{},
			Cursor:  "",
		}, nil
	}

	budgets := make([]model.Budget, len(dbBudgets))
	for i, dbBudget := range dbBudgets {
		budgets[i] = model.Budget{
			ID:        dbBudget.ID,
			CreatedAt: dbBudget.CreatedAt,
			UpdatedAt: dbBudget.UpdatedAt,

			Amount:     decimal.RequireFromString(dbBudget.Amount),
			Goal:       decimal.RequireFromString(dbBudget.Goal),
			StartDate:  dbBudget.StartDate,
			EndDate:    dbBudget.EndDate,
			UserID:     dbBudget.UserID,
			CategoryID: dbBudget.CategoryID,
		}
	}

	cursor := ""
	if len(budgets) == int(page.Limit) {
		cursor = encodeCursor(budgets[len(budgets)-1].CreatedAt, budgets[len(budgets)-1].ID)
	}

	return FindResult{
		Budgets: budgets,
		Cursor:  cursor,
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
