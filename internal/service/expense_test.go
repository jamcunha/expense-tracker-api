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

func TestExpenseService(t *testing.T) {
	var (
		mockDB         *mocks.MockDB
		mockQuerier    *mocks.MockQuerier
		expenseService service.Expense
	)

	now := time.Now()
	expectedExpense := repository.Expense{
		ID:          uuid.New(),
		CreatedAt:   now,
		UpdatedAt:   now,
		Description: "test expense",
		Amount:      decimal.NewFromFloat(1.23),
		CategoryID:  uuid.New(),
		UserID:      uuid.New(),
	}

	setup := func() {
		mockDB = &mocks.MockDB{}
		mockQuerier = &mocks.MockQuerier{}

		expenseService = service.Expense{
			DB:      mockDB,
			Queries: mockQuerier,
		}
	}

	t.Run("it should create an expense", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("CreateExpense", ctx, mock.Anything).Return(expectedExpense, nil)

		e, err := expenseService.Create(
			ctx,
			expectedExpense.UserID,
			expectedExpense.Description,
			expectedExpense.Amount,
			expectedExpense.CategoryID,
		)

		assert.Nil(t, err)
		assertExpenseEqual(t, expectedExpense, e)
	})

	t.Run("it should return an expense given it's id", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetExpenseByID", ctx, mock.Anything).Return(expectedExpense, nil)

		e, err := expenseService.GetByID(ctx, expectedExpense.ID, expectedExpense.UserID)

		assert.Nil(t, err)
		assertExpenseEqual(t, expectedExpense, e)
	})

	t.Run("it should return an error if expense is not found", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetExpenseByID", ctx, mock.Anything).
			Return(repository.Expense{}, pgx.ErrNoRows)

		e, err := expenseService.GetByID(ctx, expectedExpense.ID, expectedExpense.UserID)

		assert.ErrorIs(t, err, service.ErrExpenseNotFound)
		assert.Empty(t, e)
	})

	t.Run("it should return all expenses belonging to a user", func(t *testing.T) {
		setup()

		ctx := context.Background()

		userID := uuid.New()
		now := time.Now()

		expectedExpenses := []repository.Expense{
			{
				ID:          uuid.New(),
				CreatedAt:   now,
				UpdatedAt:   now,
				Description: "Bread",
				Amount:      decimal.NewFromFloat(0.32),
				CategoryID:  uuid.New(),
				UserID:      userID,
			},
			{
				ID:          uuid.New(),
				CreatedAt:   now.Add(5 * time.Minute),
				UpdatedAt:   now.Add(5 * time.Minute),
				Description: "Bolt",
				Amount:      decimal.NewFromFloat(1.45),
				CategoryID:  uuid.New(),
				UserID:      userID,
			},
			{
				ID:          uuid.New(),
				CreatedAt:   now.Add(10 * time.Minute),
				UpdatedAt:   now.Add(10 * time.Minute),
				Description: "Subscription",
				Amount:      decimal.NewFromFloat(13.25),
				CategoryID:  uuid.New(),
				UserID:      userID,
			},
		}

		mockQuerier.On("GetUserExpenses", ctx, mock.Anything).Return(expectedExpenses, nil)

		expenses, err := expenseService.GetAll(ctx, userID, 10, "")

		assert.Nil(t, err)
		assert.Len(t, expenses, 3)

		for i := range expenses {
			assertExpenseEqual(t, expectedExpenses[i], expenses[i])
		}
	})

	t.Run("it should return all expenses belonging to a user paged", func(t *testing.T) {
		setup()

		ctx := context.Background()

		userID := uuid.New()
		now := time.Now()

		expectedExpenses := []repository.Expense{
			{
				ID:          uuid.New(),
				CreatedAt:   now,
				UpdatedAt:   now,
				Description: "Bread",
				Amount:      decimal.NewFromFloat(0.32),
				CategoryID:  uuid.New(),
				UserID:      userID,
			},
			{
				ID:          uuid.New(),
				CreatedAt:   now.Add(5 * time.Minute),
				UpdatedAt:   now.Add(5 * time.Minute),
				Description: "Bolt",
				Amount:      decimal.NewFromFloat(1.45),
				CategoryID:  uuid.New(),
				UserID:      userID,
			},
			{
				ID:          uuid.New(),
				CreatedAt:   now.Add(10 * time.Minute),
				UpdatedAt:   now.Add(10 * time.Minute),
				Description: "Subscription",
				Amount:      decimal.NewFromFloat(13.25),
				CategoryID:  uuid.New(),
				UserID:      userID,
			},
		}

		mockCursor := "MjAyNS0wMS0xN1QwNjoxMTo1MC43OTEzNDU4MzJaLDk4OWU4ZjU4LTRlMWMtNDMxZS05MTI0LWJlNDU3ZDIwMmU5OQ=="
		mockQuerier.On("GetUserExpensesPaged", ctx, mock.Anything).Return(expectedExpenses[2:], nil)

		expenses, err := expenseService.GetAll(ctx, userID, 2, mockCursor)

		assert.Nil(t, err)
		assert.Len(t, expenses, 1)

		for i := range expenses {
			assertExpenseEqual(t, expectedExpenses[2+i], expenses[i])
		}
	})

	t.Run("it should return an error when user has no expenses", func(t *testing.T) {
		setup()

		ctx := context.Background()

		userID := uuid.New()

		mockQuerier.On("GetUserExpenses", ctx, mock.Anything).
			Return([]repository.Expense{}, pgx.ErrNoRows)

		expenses, err := expenseService.GetAll(ctx, userID, 10, "")

		assert.ErrorIs(t, err, service.ErrExpenseNotFound)
		assert.Empty(t, expenses)
	})

	t.Run("it should return all expenses from a category", func(t *testing.T) {
		setup()

		ctx := context.Background()

		userID := uuid.New()
		categoryID := uuid.New()
		now := time.Now()

		expectedExpenses := []repository.Expense{
			{
				ID:          uuid.New(),
				CreatedAt:   now,
				UpdatedAt:   now,
				Description: "Bread",
				Amount:      decimal.NewFromFloat(0.32),
				CategoryID:  categoryID,
				UserID:      userID,
			},
			{
				ID:          uuid.New(),
				CreatedAt:   now.Add(5 * time.Minute),
				UpdatedAt:   now.Add(5 * time.Minute),
				Description: "Chocolate",
				Amount:      decimal.NewFromFloat(1.45),
				CategoryID:  categoryID,
				UserID:      userID,
			},
			{
				ID:          uuid.New(),
				CreatedAt:   now.Add(10 * time.Minute),
				UpdatedAt:   now.Add(10 * time.Minute),
				Description: "Meat",
				Amount:      decimal.NewFromFloat(13.25),
				CategoryID:  categoryID,
				UserID:      userID,
			},
		}

		mockQuerier.On("GetCategoryExpenses", ctx, mock.Anything).Return(expectedExpenses, nil)

		expenses, err := expenseService.GetByCategory(ctx, categoryID, userID, 10, "")

		assert.Nil(t, err)
		assert.Len(t, expenses, 3)

		for i := range expenses {
			assertExpenseEqual(t, expectedExpenses[i], expenses[i])
		}
	})

	t.Run("it should return all expenses from a category, paged", func(t *testing.T) {
		setup()

		ctx := context.Background()

		userID := uuid.New()
		categoryID := uuid.New()
		now := time.Now()

		expectedExpenses := []repository.Expense{
			{
				ID:          uuid.New(),
				CreatedAt:   now,
				UpdatedAt:   now,
				Description: "Bread",
				Amount:      decimal.NewFromFloat(0.32),
				CategoryID:  categoryID,
				UserID:      userID,
			},
			{
				ID:          uuid.New(),
				CreatedAt:   now.Add(5 * time.Minute),
				UpdatedAt:   now.Add(5 * time.Minute),
				Description: "Chocolate",
				Amount:      decimal.NewFromFloat(1.45),
				CategoryID:  categoryID,
				UserID:      userID,
			},
			{
				ID:          uuid.New(),
				CreatedAt:   now.Add(10 * time.Minute),
				UpdatedAt:   now.Add(10 * time.Minute),
				Description: "Meat",
				Amount:      decimal.NewFromFloat(13.25),
				CategoryID:  categoryID,
				UserID:      userID,
			},
		}

		mockCursor := "MjAyNS0wMS0xN1QwNjoxMTo1MC43OTEzNDU4MzJaLDk4OWU4ZjU4LTRlMWMtNDMxZS05MTI0LWJlNDU3ZDIwMmU5OQ=="
		mockQuerier.On("GetCategoryExpensesPaged", ctx, mock.Anything).
			Return(expectedExpenses[2:], nil)

		expenses, err := expenseService.GetByCategory(ctx, categoryID, userID, 2, mockCursor)

		assert.Nil(t, err)
		assert.Len(t, expenses, 1)

		for i := range expenses {
			assertExpenseEqual(t, expectedExpenses[2+i], expenses[i])
		}
	})

	t.Run("it should return an error when there are no expenses in category", func(t *testing.T) {
		setup()

		ctx := context.Background()

		userID := uuid.New()
		categoryID := uuid.New()

		mockQuerier.On("GetCategoryExpenses", ctx, mock.Anything).
			Return([]repository.Expense{}, pgx.ErrNoRows)

		expenses, err := expenseService.GetByCategory(ctx, categoryID, userID, 10, "")

		assert.ErrorIs(t, err, service.ErrExpenseNotFound)
		assert.Empty(t, expenses)
	})

	t.Run("it should update the description of the expense", func(t *testing.T) {
		setup()

		ctx := context.Background()

		updatedExpense := expectedExpense
		updatedExpense.UpdatedAt = time.Now()
		updatedExpense.Description = "update"

		mockTx := mocks.MockTx{}
		mockTx.On("Rollback", ctx).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)

		mockDB.On("Begin", ctx).Return(&mockTx, nil)

		mockQuerier.On("WithTx", mock.Anything).Return(mockQuerier)
		mockQuerier.On("GetExpenseByID", ctx, mock.Anything).Return(expectedExpense, nil)
		mockQuerier.On("UpdateExpense", ctx, mock.Anything).Return(updatedExpense, nil)
		mockQuerier.On("UpdateBudgetAmount", ctx, mock.Anything).Return(nil)

		e, err := expenseService.Update(
			ctx,
			expectedExpense.ID,
			uuid.Nil,
			expectedExpense.UserID,
			updatedExpense.Description,
			decimal.Zero,
		)

		assert.Nil(t, err)
		assert.True(t, e.UpdatedAt.After(expectedExpense.UpdatedAt))

		assertExpenseEqual(t, updatedExpense, e)
	})

	t.Run("it should update the amount of the expense", func(t *testing.T) {
		setup()

		ctx := context.Background()

		updatedExpense := expectedExpense
		updatedExpense.UpdatedAt = time.Now()
		updatedExpense.Amount = decimal.NewFromFloat(32.40)

		mockTx := mocks.MockTx{}
		mockTx.On("Rollback", ctx).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)

		mockDB.On("Begin", ctx).Return(&mockTx, nil)

		mockQuerier.On("WithTx", mock.Anything).Return(mockQuerier)
		mockQuerier.On("GetExpenseByID", ctx, mock.Anything).Return(expectedExpense, nil)
		mockQuerier.On("UpdateExpense", ctx, mock.Anything).Return(updatedExpense, nil)
		mockQuerier.On("UpdateBudgetAmount", ctx, mock.Anything).Return(nil)

		e, err := expenseService.Update(
			ctx,
			expectedExpense.ID,
			uuid.Nil,
			expectedExpense.UserID,
			"",
			updatedExpense.Amount,
		)

		assert.Nil(t, err)
		assert.True(t, e.UpdatedAt.After(expectedExpense.UpdatedAt))

		assertExpenseEqual(t, updatedExpense, e)
	})

	t.Run("it should update the category of the expense", func(t *testing.T) {
		setup()

		ctx := context.Background()

		updatedExpense := expectedExpense
		updatedExpense.UpdatedAt = time.Now()
		updatedExpense.CategoryID = uuid.New()

		mockTx := mocks.MockTx{}
		mockTx.On("Rollback", ctx).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)

		mockDB.On("Begin", ctx).Return(&mockTx, nil)

		mockQuerier.On("WithTx", mock.Anything).Return(mockQuerier)
		mockQuerier.On("GetExpenseByID", ctx, mock.Anything).Return(expectedExpense, nil)
		mockQuerier.On("UpdateExpense", ctx, mock.Anything).Return(updatedExpense, nil)
		mockQuerier.On("UpdateBudgetAmount", ctx, mock.Anything).Return(nil)

		e, err := expenseService.Update(
			ctx,
			expectedExpense.ID,
			updatedExpense.CategoryID,
			expectedExpense.UserID,
			"",
			decimal.Zero,
		)

		assert.Nil(t, err)
		assert.True(t, e.UpdatedAt.After(expectedExpense.UpdatedAt))

		assertExpenseEqual(t, updatedExpense, e)
	})

	t.Run("it should return an error when there is no expense to update", func(t *testing.T) {
		setup()

		// There are more error paths to test

		ctx := context.Background()

		mockTx := mocks.MockTx{}
		mockTx.On("Rollback", ctx).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)

		mockDB.On("Begin", ctx).Return(&mockTx, nil)

		mockQuerier.On("WithTx", mock.Anything).Return(mockQuerier)
		mockQuerier.On("GetExpenseByID", ctx, mock.Anything).
			Return(repository.Expense{}, pgx.ErrNoRows)

		e, err := expenseService.Update(
			ctx,
			expectedExpense.ID,
			uuid.Nil,
			expectedExpense.UserID,
			"update",
			decimal.Zero,
		)

		assert.ErrorIs(t, err, service.ErrExpenseNotFound)
		assert.Empty(t, e)
	})

	t.Run("it should delete an expense", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockTx := mocks.MockTx{}
		mockTx.On("Rollback", ctx).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)

		mockDB.On("Begin", ctx).Return(&mockTx, nil)

		mockQuerier.On("WithTx", mock.Anything).Return(mockQuerier)
		mockQuerier.On("DeleteExpense", ctx, mock.Anything).
			Return(expectedExpense, nil)
		mockQuerier.On("UpdateBudgetAmount", ctx, mock.Anything).
			Return(nil)

		e, err := expenseService.DeleteByID(ctx, expectedExpense.ID, expectedExpense.UserID)

		assert.Nil(t, err)

		assertExpenseEqual(t, expectedExpense, e)
	})

	t.Run("it should return an error when there is no expense to delete", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockTx := mocks.MockTx{}
		mockTx.On("Rollback", ctx).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)

		mockDB.On("Begin", ctx).Return(&mockTx, nil)

		mockQuerier.On("WithTx", mock.Anything).Return(mockQuerier)
		mockQuerier.On("DeleteExpense", ctx, mock.Anything).
			Return(repository.Expense{}, pgx.ErrNoRows)

		e, err := expenseService.DeleteByID(ctx, expectedExpense.ID, expectedExpense.UserID)

		assert.ErrorIs(t, err, service.ErrExpenseNotFound)
		assert.Empty(t, e)
	})
}

func assertExpenseEqual(t *testing.T, expected repository.Expense, actual repository.Expense) {
	t.Helper()

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.Amount, actual.Amount)
	assert.Equal(t, expected.CategoryID, actual.CategoryID)
	assert.Equal(t, expected.UserID, actual.UserID)
}
