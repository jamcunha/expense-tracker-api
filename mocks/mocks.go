package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jamcunha/expense-tracker/internal"
	"github.com/jamcunha/expense-tracker/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

// TODO: extend mockDb (or mockQuerier) to store data

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockDB) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) WithTx(tx pgx.Tx) internal.Querier {
	args := m.Called(tx)
	return args.Get(0).(internal.Querier)
}

func (m *MockQuerier) CreateBudget(
	ctx context.Context,
	arg repository.CreateBudgetParams,
) (repository.Budget, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Budget), args.Error(1)
}

func (m *MockQuerier) CreateCategory(
	ctx context.Context,
	arg repository.CreateCategoryParams,
) (repository.Category, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Category), args.Error(1)
}

func (m *MockQuerier) CreateExpense(
	ctx context.Context,
	arg repository.CreateExpenseParams,
) (repository.Expense, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Expense), args.Error(1)
}

func (m *MockQuerier) CreateUser(
	ctx context.Context,
	arg repository.CreateUserParams,
) (repository.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockQuerier) DeleteBudget(
	ctx context.Context,
	arg repository.DeleteBudgetParams,
) (repository.Budget, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Budget), args.Error(1)
}

func (m *MockQuerier) DeleteCategory(
	ctx context.Context,
	arg repository.DeleteCategoryParams,
) (repository.Category, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Category), args.Error(1)
}

func (m *MockQuerier) DeleteExpense(
	ctx context.Context,
	arg repository.DeleteExpenseParams,
) (repository.Expense, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Expense), args.Error(1)
}

func (m *MockQuerier) DeleteUser(ctx context.Context, id uuid.UUID) (repository.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockQuerier) GetBudgetByID(
	ctx context.Context,
	arg repository.GetBudgetByIDParams,
) (repository.Budget, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Budget), args.Error(1)
}

func (m *MockQuerier) GetCategoryByID(
	ctx context.Context,
	arg repository.GetCategoryByIDParams,
) (repository.Category, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Category), args.Error(1)
}

func (m *MockQuerier) GetCategoryExpenses(
	ctx context.Context,
	arg repository.GetCategoryExpensesParams,
) ([]repository.Expense, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repository.Expense), args.Error(1)
}

func (m *MockQuerier) GetCategoryExpensesPaged(
	ctx context.Context,
	arg repository.GetCategoryExpensesPagedParams,
) ([]repository.Expense, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repository.Expense), args.Error(1)
}

func (m *MockQuerier) GetExpenseByID(
	ctx context.Context,
	arg repository.GetExpenseByIDParams,
) (repository.Expense, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Expense), args.Error(1)
}

func (m *MockQuerier) GetTotalSpent(
	ctx context.Context,
	arg repository.GetTotalSpentParams,
) (decimal.Decimal, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockQuerier) GetTotalSpentInCategory(
	ctx context.Context,
	arg repository.GetTotalSpentInCategoryParams,
) (decimal.Decimal, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockQuerier) GetUserBudgets(
	ctx context.Context,
	arg repository.GetUserBudgetsParams,
) ([]repository.Budget, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repository.Budget), args.Error(1)
}

func (m *MockQuerier) GetUserBudgetsPaged(
	ctx context.Context,
	arg repository.GetUserBudgetsPagedParams,
) ([]repository.Budget, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repository.Budget), args.Error(1)
}

func (m *MockQuerier) GetUserByEmail(ctx context.Context, email string) (repository.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockQuerier) GetUserByID(ctx context.Context, id uuid.UUID) (repository.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockQuerier) GetUserCategories(
	ctx context.Context,
	arg repository.GetUserCategoriesParams,
) ([]repository.Category, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repository.Category), args.Error(1)
}

func (m *MockQuerier) GetUserCategoriesPaged(
	ctx context.Context,
	arg repository.GetUserCategoriesPagedParams,
) ([]repository.Category, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repository.Category), args.Error(1)
}

func (m *MockQuerier) GetUserExpenses(
	ctx context.Context,
	arg repository.GetUserExpensesParams,
) ([]repository.Expense, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repository.Expense), args.Error(1)
}

func (m *MockQuerier) GetUserExpensesPaged(
	ctx context.Context,
	arg repository.GetUserExpensesPagedParams,
) ([]repository.Expense, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repository.Expense), args.Error(1)
}

func (m *MockQuerier) UpdateBudgetAmount(
	ctx context.Context,
	arg repository.UpdateBudgetAmountParams,
) error {
	args := m.Called(ctx, arg)
	return args.Error(1)
}

func (m *MockQuerier) UpdateCategory(
	ctx context.Context,
	arg repository.UpdateCategoryParams,
) (repository.Category, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Category), args.Error(1)
}

func (m *MockQuerier) UpdateExpense(
	ctx context.Context,
	arg repository.UpdateExpenseParams,
) (repository.Expense, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Expense), args.Error(1)
}
