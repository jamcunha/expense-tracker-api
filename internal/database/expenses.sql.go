// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: expenses.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createExpense = `-- name: CreateExpense :one

INSERT INTO expenses (id, created_at, updated_at, description, amount, category_id, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at, description, amount, category_id, user_id
`

type CreateExpenseParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
	Amount      string
	CategoryID  uuid.UUID
	UserID      uuid.UUID
}

// TODO: add query (and route) to get all expenses in a time interval
func (q *Queries) CreateExpense(ctx context.Context, arg CreateExpenseParams) (Expense, error) {
	row := q.db.QueryRowContext(ctx, createExpense,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Description,
		arg.Amount,
		arg.CategoryID,
		arg.UserID,
	)
	var i Expense
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
		&i.Amount,
		&i.CategoryID,
		&i.UserID,
	)
	return i, err
}

const deleteExpense = `-- name: DeleteExpense :one
DELETE FROM expenses WHERE id = $1 AND user_id = $2 RETURNING id, created_at, updated_at, description, amount, category_id, user_id
`

type DeleteExpenseParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteExpense(ctx context.Context, arg DeleteExpenseParams) (Expense, error) {
	row := q.db.QueryRowContext(ctx, deleteExpense, arg.ID, arg.UserID)
	var i Expense
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
		&i.Amount,
		&i.CategoryID,
		&i.UserID,
	)
	return i, err
}

const getCategoryExpenses = `-- name: GetCategoryExpenses :many
SELECT id, created_at, updated_at, description, amount, category_id, user_id FROM expenses WHERE category_id = $1 AND user_id = $2
ORDER BY created_at DESC, id DESC
LIMIT $3
`

type GetCategoryExpensesParams struct {
	CategoryID uuid.UUID
	UserID     uuid.UUID
	Limit      int32
}

func (q *Queries) GetCategoryExpenses(ctx context.Context, arg GetCategoryExpensesParams) ([]Expense, error) {
	rows, err := q.db.QueryContext(ctx, getCategoryExpenses, arg.CategoryID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Expense
	for rows.Next() {
		var i Expense
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Description,
			&i.Amount,
			&i.CategoryID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCategoryExpensesPaged = `-- name: GetCategoryExpensesPaged :many
SELECT id, created_at, updated_at, description, amount, category_id, user_id FROM expenses WHERE category_id = $1 AND user_id = $2
AND created_at <= $3 AND id < $4
ORDER BY created_at DESC, id DESC
LIMIT $5
`

type GetCategoryExpensesPagedParams struct {
	CategoryID uuid.UUID
	UserID     uuid.UUID
	CreatedAt  time.Time
	ID         uuid.UUID
	Limit      int32
}

func (q *Queries) GetCategoryExpensesPaged(ctx context.Context, arg GetCategoryExpensesPagedParams) ([]Expense, error) {
	rows, err := q.db.QueryContext(ctx, getCategoryExpensesPaged,
		arg.CategoryID,
		arg.UserID,
		arg.CreatedAt,
		arg.ID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Expense
	for rows.Next() {
		var i Expense
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Description,
			&i.Amount,
			&i.CategoryID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getExpenseByID = `-- name: GetExpenseByID :one
SELECT id, created_at, updated_at, description, amount, category_id, user_id FROM expenses WHERE id = $1 AND user_id = $2
`

type GetExpenseByIDParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) GetExpenseByID(ctx context.Context, arg GetExpenseByIDParams) (Expense, error) {
	row := q.db.QueryRowContext(ctx, getExpenseByID, arg.ID, arg.UserID)
	var i Expense
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
		&i.Amount,
		&i.CategoryID,
		&i.UserID,
	)
	return i, err
}

const getTotalSpent = `-- name: GetTotalSpent :one

SELECT CAST(SUM(amount) AS NUMERIC(10, 4)) FROM expenses
WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3
`

type GetTotalSpentParams struct {
	UserID      uuid.UUID
	CreatedAt   time.Time
	CreatedAt_2 time.Time
}

// NOTE: both folowing queries are private to the API, not used by the client (might be made public in the future)
func (q *Queries) GetTotalSpent(ctx context.Context, arg GetTotalSpentParams) (string, error) {
	row := q.db.QueryRowContext(ctx, getTotalSpent, arg.UserID, arg.CreatedAt, arg.CreatedAt_2)
	var column_1 string
	err := row.Scan(&column_1)
	return column_1, err
}

const getTotalSpentInCategory = `-- name: GetTotalSpentInCategory :one
SELECT CAST(SUM(amount) AS NUMERIC(10, 4)) FROM expenses
WHERE user_id = $1 AND category_id = $2 AND created_at >= $3 AND created_at <= $4
`

type GetTotalSpentInCategoryParams struct {
	UserID      uuid.UUID
	CategoryID  uuid.UUID
	CreatedAt   time.Time
	CreatedAt_2 time.Time
}

func (q *Queries) GetTotalSpentInCategory(ctx context.Context, arg GetTotalSpentInCategoryParams) (string, error) {
	row := q.db.QueryRowContext(ctx, getTotalSpentInCategory,
		arg.UserID,
		arg.CategoryID,
		arg.CreatedAt,
		arg.CreatedAt_2,
	)
	var column_1 string
	err := row.Scan(&column_1)
	return column_1, err
}

const getUserExpenses = `-- name: GetUserExpenses :many
SELECT id, created_at, updated_at, description, amount, category_id, user_id FROM expenses WHERE user_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2
`

type GetUserExpensesParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetUserExpenses(ctx context.Context, arg GetUserExpensesParams) ([]Expense, error) {
	rows, err := q.db.QueryContext(ctx, getUserExpenses, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Expense
	for rows.Next() {
		var i Expense
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Description,
			&i.Amount,
			&i.CategoryID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserExpensesPaged = `-- name: GetUserExpensesPaged :many
SELECT id, created_at, updated_at, description, amount, category_id, user_id FROM expenses WHERE user_id = $1
AND created_at <= $2 AND id < $3
ORDER BY created_at DESC, id DESC
LIMIT $4
`

type GetUserExpensesPagedParams struct {
	UserID    uuid.UUID
	CreatedAt time.Time
	ID        uuid.UUID
	Limit     int32
}

func (q *Queries) GetUserExpensesPaged(ctx context.Context, arg GetUserExpensesPagedParams) ([]Expense, error) {
	rows, err := q.db.QueryContext(ctx, getUserExpensesPaged,
		arg.UserID,
		arg.CreatedAt,
		arg.ID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Expense
	for rows.Next() {
		var i Expense
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Description,
			&i.Amount,
			&i.CategoryID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateExpense = `-- name: UpdateExpense :one
UPDATE expenses SET description = $1, amount = $2, category_id = $3, updated_at = $4
WHERE id = $5 AND user_id = $6 RETURNING id, created_at, updated_at, description, amount, category_id, user_id
`

type UpdateExpenseParams struct {
	Description string
	Amount      string
	CategoryID  uuid.UUID
	UpdatedAt   time.Time
	ID          uuid.UUID
	UserID      uuid.UUID
}

func (q *Queries) UpdateExpense(ctx context.Context, arg UpdateExpenseParams) (Expense, error) {
	row := q.db.QueryRowContext(ctx, updateExpense,
		arg.Description,
		arg.Amount,
		arg.CategoryID,
		arg.UpdatedAt,
		arg.ID,
		arg.UserID,
	)
	var i Expense
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
		&i.Amount,
		&i.CategoryID,
		&i.UserID,
	)
	return i, err
}
