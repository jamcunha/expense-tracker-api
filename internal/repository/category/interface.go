package category

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/model"
)

type Repo interface {
	Create(ctx context.Context, category model.Category) (model.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (model.Category, error)
	FindAll(ctx context.Context, userID uuid.UUID, page FindAllPage) (FindResult, error)
}

// TODO: Not sure about adding cursor encoding and decoding to interface
// encodeCursor(string) (time.Time, uuid.UUID, error)
// decodeCursor(time.Time, uuid.UUID) string

type FindAllPage struct {
	Limit  int32
	Cursor string
}

type FindResult struct {
	Categories []model.Category
	Cursor     string
}

var ErrNotFound = errors.New("category does not exist")

var ErrInvalidCursor = errors.New("invalid cursor")
