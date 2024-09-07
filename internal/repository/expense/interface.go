package expense

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/model"
)

type Repo interface {
	Create(ctx context.Context, expense model.Expense) (model.Expense, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (model.Expense, error)
	FindByCategory(ctx context.Context, categoryID uuid.UUID, page FindAllPage) (FindResult, error)
	FindAll(ctx context.Context, userID uuid.UUID, page FindAllPage) (FindResult, error)
}

type FindAllPage struct {
	Limit  int32
	Cursor string
}

type FindResult struct {
	Expenses []model.Expense
	Cursor   string
}

var ErrNotFound = errors.New("expense does not exist")

var ErrInvalidCursor = errors.New("invalid cursor")
