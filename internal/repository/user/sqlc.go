package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/database"
	"github.com/jamcunha/expense-tracker/internal/model"
)

// TODO: Implement update

type SqlcRepo struct {
	DB *database.Queries
}

func (s *SqlcRepo) Create(ctx context.Context, user model.User) error {
	_, err := s.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
	})

	return err
}

func (s *SqlcRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return s.DB.DeleteUser(ctx, id)
}

func (s *SqlcRepo) FindByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	dbUser, err := s.DB.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrNotFound
	} else if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
	}, nil
}

func (s *SqlcRepo) FindByEmail(ctx context.Context, email string) (model.User, error) {
	dbUser, err := s.DB.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrNotFound
	} else if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
	}, nil
}