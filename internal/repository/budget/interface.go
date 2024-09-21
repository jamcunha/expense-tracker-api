package budget

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/model"
)

type Repo interface {
	Create(ctx context.Context, budget model.Budget) (model.Budget, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Budget, error)
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Budget, error)
	FindAll(ctx context.Context, userID uuid.UUID, page FindAllPage) (FindResult, error)
}

type FindAllPage struct {
	Limit  int32
	Cursor string
}

type FindResult struct {
	Budgets []model.Budget
	Cursor  string
}

var ErrNotFound = errors.New("budget does not exist")

var ErrInvalidCursor = errors.New("invalid cursor")
