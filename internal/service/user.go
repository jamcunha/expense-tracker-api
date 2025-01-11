package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jamcunha/expense-tracker/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	DB      *pgx.Conn
	Queries *repository.Queries
}

func (s *User) GetByID(ctx context.Context, id uuid.UUID) (repository.User, error) {
	u, err := s.Queries.GetUserByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.User{}, ErrUserNotFound
	} else if err != nil {
		return repository.User{}, err
	}

	return u, nil
}

func (s *User) Create(ctx context.Context, name, email, password string) (repository.User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return repository.User{}, err
	}

	now := time.Now()
	u, err := s.Queries.CreateUser(ctx, repository.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Email:     email,
		Password:  string(encryptedPassword),
	})
	if err != nil {
		return repository.User{}, err
	}

	return u, nil
}

func (s *User) DeleteByID(ctx context.Context, id uuid.UUID) (repository.User, error) {
	u, err := s.Queries.DeleteUser(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.User{}, ErrUserNotFound
	}
	if err != nil {
		return repository.User{}, err
	}

	return u, nil
}
