package expense

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/model"
	"github.com/shopspring/decimal"
)

type Repo interface {
	Create(ctx context.Context, expense model.Expense) (model.Expense, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Expense, error)
	Update(ctx context.Context, update UpdateExpense) (model.Expense, error)
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Expense, error)
	FindByCategory(
		ctx context.Context,
		categoryID uuid.UUID,
		userID uuid.UUID,
		page FindAllPage,
	) (FindResult, error)
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

type UpdateExpense struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Description string
	Amount      decimal.Decimal
	CategoryID  uuid.UUID
}

var ErrNotFound = errors.New("expense does not exist")

var ErrInvalidCursor = errors.New("invalid cursor")
