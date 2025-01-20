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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBudgetService(t *testing.T) {
	var (
		mockDB        *mocks.MockDB
		mockQuerier   *mocks.MockQuerier
		budgetService service.Budget
	)

	now := time.Now()
	expectedBudget := repository.Budget{
		ID:         uuid.New(),
		CreatedAt:  now,
		UpdatedAt:  now,
		Amount:     decimal.NewFromFloat(150.50),
		StartDate:  now,
		EndDate:    now.Add(48 * time.Hour),
		UserID:     uuid.New(),
		CategoryID: uuid.New(),
	}

	setup := func() {
		mockDB = &mocks.MockDB{}
		mockQuerier = &mocks.MockQuerier{}

		budgetService = service.Budget{
			DB:      mockDB,
			Queries: mockQuerier,
		}
	}

	t.Run("it should create a new budget", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockTx := mocks.MockTx{}
		mockTx.On("Rollback", ctx).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)

		mockDB.On("Begin", ctx).Return(&mockTx, nil)

		mockQuerier.On("WithTx", mock.Anything).Return(mockQuerier)
		mockQuerier.On("GetTotalSpentInCategory", ctx, mock.Anything).
			Return(decimal.NewFromFloat(16.4), nil)
		mockQuerier.On("CreateBudget", ctx, mock.Anything).Return(expectedBudget, nil)

		b, err := budgetService.Create(
			ctx,
			expectedBudget.UserID,
			expectedBudget.CategoryID,
			expectedBudget.Amount,
			expectedBudget.StartDate,
			expectedBudget.EndDate,
		)

		assert.Nil(t, err)
		assertBudgetEqual(t, expectedBudget, b)
	})

	t.Run("it should fail to create a budget if category does not exist", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockTx := mocks.MockTx{}
		mockTx.On("Rollback", ctx).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)

		mockDB.On("Begin", ctx).Return(&mockTx, nil)

		mockQuerier.On("WithTx", mock.Anything).Return(mockQuerier)
		mockQuerier.On("GetTotalSpentInCategory", ctx, mock.Anything).
			Return(decimal.Zero, pgx.ErrNoRows)

		b, err := budgetService.Create(
			ctx,
			expectedBudget.UserID,
			expectedBudget.CategoryID,
			expectedBudget.Amount,
			expectedBudget.StartDate,
			expectedBudget.EndDate,
		)

		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
		assert.Empty(t, b)
	})

	t.Run("it should get a budget given it's id", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetBudgetByID", ctx, mock.Anything).Return(expectedBudget, nil)

		b, err := budgetService.GetByID(ctx, expectedBudget.ID, expectedBudget.UserID)

		assert.Nil(t, err)
		assertBudgetEqual(t, expectedBudget, b)
	})

	t.Run("it should return an error if budget does not exist", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetBudgetByID", ctx, mock.Anything).
			Return(repository.Budget{}, pgx.ErrNoRows)

		b, err := budgetService.GetByID(ctx, expectedBudget.ID, expectedBudget.UserID)

		assert.ErrorIs(t, err, service.ErrBudgetNotFound)
		assert.Empty(t, b)
	})

	t.Run("it should return all budgets belonging to an user", func(t *testing.T) {
		setup()

		ctx := context.Background()

		userID := uuid.New()
		expectedBudgets := []repository.Budget{
			{
				ID:         uuid.New(),
				CreatedAt:  now,
				UpdatedAt:  now,
				Amount:     decimal.NewFromFloat(150.50),
				StartDate:  now,
				EndDate:    now.Add(48 * time.Hour),
				UserID:     userID,
				CategoryID: uuid.New(),
			},
			{
				ID:         uuid.New(),
				CreatedAt:  now.Add(5 * time.Minute),
				UpdatedAt:  now.Add(5 * time.Minute),
				Amount:     decimal.NewFromFloat(750.50),
				StartDate:  now,
				EndDate:    now.Add(72 * time.Hour),
				UserID:     userID,
				CategoryID: uuid.New(),
			},
			{
				ID:         uuid.New(),
				CreatedAt:  now.Add(10 * time.Minute),
				UpdatedAt:  now.Add(10 * time.Minute),
				Amount:     decimal.NewFromFloat(450.20),
				StartDate:  now,
				EndDate:    now.Add(720 * time.Hour),
				UserID:     userID,
				CategoryID: uuid.New(),
			},
		}

		mockQuerier.On("GetUserBudgets", ctx, mock.Anything).Return(expectedBudgets, nil)

		budgets, err := budgetService.GetAll(ctx, userID, 10, "")

		assert.Nil(t, err)
		assert.Len(t, budgets, 3)

		for i := range expectedBudgets {
			assertBudgetEqual(t, expectedBudgets[i], budgets[i])
		}
	})

	t.Run("it should return all budgets belonging to an user, paged", func(t *testing.T) {
		setup()

		ctx := context.Background()

		userID := uuid.New()
		expectedBudgets := []repository.Budget{
			{
				ID:         uuid.New(),
				CreatedAt:  now.Add(10 * time.Minute),
				UpdatedAt:  now.Add(10 * time.Minute),
				Amount:     decimal.NewFromFloat(450.20),
				StartDate:  now,
				EndDate:    now.Add(720 * time.Hour),
				UserID:     userID,
				CategoryID: uuid.New(),
			},
		}

		mockQuerier.On("GetUserBudgetsPaged", ctx, mock.Anything).Return(expectedBudgets, nil)

		mockCursor := "MjAyNS0wMS0xN1QwNjoxMTo1MC43OTEzNDU4MzJaLDk4OWU4ZjU4LTRlMWMtNDMxZS05MTI0LWJlNDU3ZDIwMmU5OQ=="
		budgets, err := budgetService.GetAll(ctx, userID, 2, mockCursor)

		assert.Nil(t, err)
		assert.Len(t, budgets, 1)

		for i := range expectedBudgets {
			assertBudgetEqual(t, expectedBudgets[i], budgets[i])
		}
	})

	t.Run("it should an error if there are no budgets", func(t *testing.T) {
		setup()

		ctx := context.Background()

		userID := uuid.New()

		mockQuerier.On("GetUserBudgets", ctx, mock.Anything).
			Return([]repository.Budget{}, pgx.ErrNoRows)

		budgets, err := budgetService.GetAll(ctx, userID, 10, "")

		assert.ErrorIs(t, err, service.ErrBudgetNotFound)
		assert.Empty(t, budgets)
	})

	t.Run("it should delete a budget", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("DeleteBudget", ctx, mock.Anything).Return(expectedBudget, nil)

		b, err := budgetService.DeleteByID(ctx, expectedBudget.ID, expectedBudget.UserID)

		assert.Nil(t, err)
		assertBudgetEqual(t, expectedBudget, b)
	})

	t.Run("it should return an error if there is no budget to delete", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("DeleteBudget", ctx, mock.Anything).
			Return(repository.Budget{}, pgx.ErrNoRows)

		b, err := budgetService.DeleteByID(ctx, expectedBudget.ID, expectedBudget.UserID)

		assert.ErrorIs(t, err, service.ErrBudgetNotFound)
		assert.Empty(t, b)
	})
}

func assertBudgetEqual(t *testing.T, expected repository.Budget, actual repository.Budget) {
	t.Helper()

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
	assert.Equal(t, expected.Amount, actual.Amount)
	assert.Equal(t, expected.StartDate, actual.StartDate)
	assert.Equal(t, expected.EndDate, actual.EndDate)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.CategoryID, actual.CategoryID)
}
