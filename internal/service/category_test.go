package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jamcunha/expense-tracker/internal/repository"
	"github.com/jamcunha/expense-tracker/internal/service"
	"github.com/jamcunha/expense-tracker/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryService(t *testing.T) {
	var (
		mockDB          *mocks.MockDB
		mockQuerier     *mocks.MockQuerier
		categoryService service.Category
	)

	now := time.Now()
	expectedCategory := repository.Category{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "shop",
		UserID:    uuid.New(),
	}

	setup := func() {
		mockDB = &mocks.MockDB{}
		mockQuerier = &mocks.MockQuerier{}

		categoryService = service.Category{
			DB:      mockDB,
			Queries: mockQuerier,
		}
	}

	t.Run("it should create a category", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("CreateCategory", ctx, mock.Anything).Return(expectedCategory, nil)

		c, err := categoryService.Create(ctx, expectedCategory.Name, expectedCategory.UserID)

		assert.Nil(t, err)
		assertCategoryEqual(t, c, expectedCategory)
	})

	t.Run("it should get a category given it's id", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetCategoryByID", ctx, repository.GetCategoryByIDParams{
			ID:     expectedCategory.ID,
			UserID: expectedCategory.UserID,
		}).Return(expectedCategory, nil)

		c, err := categoryService.GetByID(ctx, expectedCategory.ID, expectedCategory.UserID)

		assert.Nil(t, err)
		assertCategoryEqual(t, expectedCategory, c)
	})

	t.Run("it should return an error when id is not found", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetCategoryByID", ctx, mock.Anything).
			Return(repository.Category{}, pgx.ErrNoRows)

		c, err := categoryService.GetByID(ctx, uuid.New(), uuid.New())

		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
		assert.Empty(t, c)
	})

	t.Run("it should get all categories belonging to a user", func(t *testing.T) {
		setup()

		ctx := context.Background()

		now := time.Now()
		userID := uuid.New()

		expectedCategories := []repository.Category{
			{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
				Name:      "cat1",
				UserID:    userID,
			},
			{
				ID:        uuid.New(),
				CreatedAt: now.Add(5 * time.Minute),
				UpdatedAt: now.Add(5 * time.Minute),
				Name:      "cat2",
				UserID:    userID,
			},
			{
				ID:        uuid.New(),
				CreatedAt: now.Add(10 * time.Minute),
				UpdatedAt: now.Add(10 * time.Minute),
				Name:      "cat3",
				UserID:    userID,
			},
		}

		mockQuerier.On("GetUserCategories", ctx, mock.Anything).Return(expectedCategories, nil)

		categories, err := categoryService.GetAll(ctx, userID, 10, "")

		assert.Nil(t, err)
		assert.Len(t, categories, 3)

		for i := range categories {
			assertCategoryEqual(t, expectedCategories[i], categories[i])
		}
	})

	t.Run("it should get all categories belonging to a user, paged", func(t *testing.T) {
		setup()

		ctx := context.Background()

		now := time.Now()
		userID := uuid.New()

		expectedCategories := []repository.Category{
			{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
				Name:      "cat1",
				UserID:    userID,
			},
			{
				ID:        uuid.New(),
				CreatedAt: now.Add(5 * time.Minute),
				UpdatedAt: now.Add(5 * time.Minute),
				Name:      "cat2",
				UserID:    userID,
			},
			{
				ID:        uuid.New(),
				CreatedAt: now.Add(10 * time.Minute),
				UpdatedAt: now.Add(10 * time.Minute),
				Name:      "cat3",
				UserID:    userID,
			},
		}

		mockQuerier.On("GetUserCategoriesPaged", ctx, mock.Anything).
			Return(expectedCategories[2:], nil)

		mockCursor := "MjAyNS0wMS0xN1QwNjoxMTo1MC43OTEzNDU4MzJaLDk4OWU4ZjU4LTRlMWMtNDMxZS05MTI0LWJlNDU3ZDIwMmU5OQ=="
		categories, err := categoryService.GetAll(ctx, userID, 2, mockCursor)

		assert.Nil(t, err)
		assert.Len(t, categories, 1)

		for i := range categories {
			assertCategoryEqual(t, expectedCategories[2+i], categories[i])
		}
	})

	t.Run("it should return an error when user has no categories", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetUserCategories", ctx, mock.Anything).
			Return([]repository.Category{}, pgx.ErrNoRows)

		categories, err := categoryService.GetAll(ctx, uuid.New(), 10, "")

		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
		assert.Empty(t, categories)
	})

	t.Run("it should update the category", func(t *testing.T) {
		setup()

		ctx := context.Background()

		updatedCategory := expectedCategory
		updatedCategory.UpdatedAt = time.Now()
		updatedCategory.Name = "update"

		mockQuerier.On("UpdateCategory", ctx, mock.Anything).Return(updatedCategory, nil)

		c, err := categoryService.Update(
			ctx,
			expectedCategory.ID,
			expectedCategory.UserID,
			"update",
		)

		assert.Nil(t, err)
		assert.True(t, c.UpdatedAt.After(expectedCategory.UpdatedAt))

		assertCategoryEqual(t, updatedCategory, c)
	})

	t.Run("it should return an error if the category to update does not exist", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("UpdateCategory", ctx, mock.Anything).
			Return(repository.Category{}, pgx.ErrNoRows)

		c, err := categoryService.Update(
			ctx,
			expectedCategory.ID,
			expectedCategory.UserID,
			"update",
		)

		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
		assert.Empty(t, c)
	})

	t.Run("it should delete a category given it's id", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("DeleteCategory", ctx, mock.Anything).Return(expectedCategory, nil)

		c, err := categoryService.DeleteByID(ctx, expectedCategory.ID, expectedCategory.UserID)

		assert.Nil(t, err)
		assertCategoryEqual(t, expectedCategory, c)
	})

	t.Run("it should return an error if category to delete does not exist", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("DeleteCategory", ctx, mock.Anything).
			Return(repository.Category{}, pgx.ErrNoRows)

		c, err := categoryService.DeleteByID(ctx, expectedCategory.ID, expectedCategory.UserID)

		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
		assert.Empty(t, c)
	})
}

func assertCategoryEqual(t *testing.T, expected repository.Category, actual repository.Category) {
	t.Helper()

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.UserID, actual.UserID)
}
