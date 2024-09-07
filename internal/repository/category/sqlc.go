package category

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
)

// TODO: Implement update

type SqlcRepo struct {
	DB      *sql.DB
	Queries *database.Queries
}

func (s *SqlcRepo) Create(ctx context.Context, category model.Category) (model.Category, error) {
	_, err := s.Queries.CreateCategory(ctx, database.CreateCategoryParams{
		ID:        category.ID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
		Name:      category.Name,
		UserID:    category.UserID,
	})

	return category, err
}

func (s *SqlcRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return s.Queries.DeleteCategory(ctx, id)
}

func (s *SqlcRepo) FindByID(ctx context.Context, id uuid.UUID) (model.Category, error) {
	dbCategory, err := s.Queries.GetCategoryByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Category{}, ErrNotFound
	} else if err != nil {
		return model.Category{}, err
	}

	return model.Category{
		ID:        dbCategory.ID,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
		Name:      dbCategory.Name,
		UserID:    dbCategory.UserID,
	}, nil
}

func (s *SqlcRepo) FindAll(
	ctx context.Context,
	userID uuid.UUID,
	page FindAllPage,
) (FindResult, error) {
	var dbCategories []database.Category
	var err error

	if page.Cursor == "" {
		dbCategories, err = s.Queries.GetUserCategories(ctx, database.GetUserCategoriesParams{
			UserID: userID,
			Limit:  page.Limit,
		})
	} else {
		t, id, err := decodeCursor(page.Cursor)
		if err != nil {
			return FindResult{}, err
		}

		dbCategories, err = s.Queries.GetUserCategoriesPaged(ctx, database.GetUserCategoriesPagedParams{
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
	if len(dbCategories) == 0 {
		return FindResult{
			Categories: []model.Category{},
			Cursor:     "",
		}, nil
	}

	categories := make([]model.Category, len(dbCategories))
	for i, dbCategory := range dbCategories {
		categories[i] = model.Category{
			ID:        dbCategory.ID,
			CreatedAt: dbCategory.CreatedAt,
			UpdatedAt: dbCategory.UpdatedAt,
			Name:      dbCategory.Name,
			UserID:    dbCategory.UserID,
		}
	}

	cursor := ""
	if len(categories) == int(page.Limit) {
		cursor = encodeCursor(
			dbCategories[len(dbCategories)-1].CreatedAt,
			dbCategories[len(dbCategories)-1].ID,
		)
	}

	return FindResult{
		Categories: categories,
		Cursor:     cursor,
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
