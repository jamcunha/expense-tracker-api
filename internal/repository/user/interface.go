package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/model"
)

type Repo interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	Delete(ctx context.Context, id uuid.UUID) (model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
}

var ErrNotFound = errors.New("user does not exist")
