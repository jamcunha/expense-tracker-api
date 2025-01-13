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
)

type Category struct {
	DB      internal.DBConn
	Queries internal.Querier
}

func (s *Category) GetByID(ctx context.Context, id, userID uuid.UUID) (repository.Category, error) {
	c, err := s.Queries.GetCategoryByID(ctx, repository.GetCategoryByIDParams{
		ID:     id,
		UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.Category{}, ErrCategoryNotFound
	} else if err != nil {
		fmt.Print("failed to insert:", err)
		return repository.Category{}, err
	}

	return c, nil
}

func (s *Category) GetAll(
	ctx context.Context,
	userID uuid.UUID,
	limit int32,
	cur string,
) ([]repository.Category, error) {
	var categories []repository.Category
	var err error

	if cur == "" {
		categories, err = s.Queries.GetUserCategories(ctx, repository.GetUserCategoriesParams{
			UserID: userID,
			Limit:  int32(limit),
		})
	} else {
		t, id, err := internal.DecodeCursor(cur)
		if err != nil {
			return []repository.Category{}, err
		}

		categories, err = s.Queries.GetUserCategoriesPaged(ctx, repository.GetUserCategoriesPagedParams{
			UserID:    userID,
			CreatedAt: t,
			ID:        id,
			Limit:     int32(limit),
		})
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return []repository.Category{}, ErrCategoryNotFound
	} else if err != nil {
		fmt.Println("failed to find:", err)
		return []repository.Category{}, err
	}

	return categories, nil
}

func (s *Category) Create(
	ctx context.Context,
	name string,
	userID uuid.UUID,
) (repository.Category, error) {
	now := time.Now()
	c, err := s.Queries.CreateCategory(ctx, repository.CreateCategoryParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		UserID:    userID,
	})
	if err != nil {
		fmt.Println("failed to insert:", err)
		return repository.Category{}, err
	}

	return c, nil
}

func (s *Category) Update(
	ctx context.Context,
	id, userID uuid.UUID,
	name string,
) (repository.Category, error) {
	c, err := s.Queries.UpdateCategory(ctx, repository.UpdateCategoryParams{
		Name:      name,
		UpdatedAt: time.Now(),
		ID:        id,
		UserID:    userID,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return repository.Category{}, ErrCategoryNotFound
	} else if err != nil {
		return repository.Category{}, err
	}

	return c, nil
}

func (s *Category) DeleteByID(
	ctx context.Context,
	id, userID uuid.UUID,
) (repository.Category, error) {
	c, err := s.Queries.DeleteCategory(ctx, repository.DeleteCategoryParams{
		ID:     id,
		UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.Category{}, ErrCategoryNotFound
	} else if err != nil {
		fmt.Println("failed to delete:", err)
		return repository.Category{}, err
	}

	return c, nil
}
