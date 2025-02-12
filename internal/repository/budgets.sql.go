// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: budgets.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const createBudget = `-- name: CreateBudget :one
INSERT INTO budgets (id, created_at, updated_at, amount, goal, start_date, end_date, user_id, category_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, created_at, updated_at, amount, goal, start_date, end_date, user_id, category_id
`

type CreateBudgetParams struct {
	ID         uuid.UUID       `json:"id"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	Amount     decimal.Decimal `json:"amount"`
	Goal       decimal.Decimal `json:"goal"`
	StartDate  time.Time       `json:"start_date"`
	EndDate    time.Time       `json:"end_date"`
	UserID     uuid.UUID       `json:"user_id"`
	CategoryID uuid.UUID       `json:"category_id"`
}

func (q *Queries) CreateBudget(ctx context.Context, arg CreateBudgetParams) (Budget, error) {
	row := q.db.QueryRow(ctx, createBudget,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Amount,
		arg.Goal,
		arg.StartDate,
		arg.EndDate,
		arg.UserID,
		arg.CategoryID,
	)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Amount,
		&i.Goal,
		&i.StartDate,
		&i.EndDate,
		&i.UserID,
		&i.CategoryID,
	)
	return i, err
}

const deleteBudget = `-- name: DeleteBudget :one
DELETE FROM budgets WHERE id = $1 AND user_id = $2 RETURNING id, created_at, updated_at, amount, goal, start_date, end_date, user_id, category_id
`

type DeleteBudgetParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteBudget(ctx context.Context, arg DeleteBudgetParams) (Budget, error) {
	row := q.db.QueryRow(ctx, deleteBudget, arg.ID, arg.UserID)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Amount,
		&i.Goal,
		&i.StartDate,
		&i.EndDate,
		&i.UserID,
		&i.CategoryID,
	)
	return i, err
}

const getBudgetByID = `-- name: GetBudgetByID :one
SELECT id, created_at, updated_at, amount, goal, start_date, end_date, user_id, category_id FROM budgets WHERE id = $1 AND user_id = $2
`

type GetBudgetByIDParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) GetBudgetByID(ctx context.Context, arg GetBudgetByIDParams) (Budget, error) {
	row := q.db.QueryRow(ctx, getBudgetByID, arg.ID, arg.UserID)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Amount,
		&i.Goal,
		&i.StartDate,
		&i.EndDate,
		&i.UserID,
		&i.CategoryID,
	)
	return i, err
}

const getUserBudgets = `-- name: GetUserBudgets :many
SELECT id, created_at, updated_at, amount, goal, start_date, end_date, user_id, category_id FROM budgets WHERE user_id = $1
ORDER BY created_at ASC, id DESC
LIMIT $2
`

type GetUserBudgetsParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
}

func (q *Queries) GetUserBudgets(ctx context.Context, arg GetUserBudgetsParams) ([]Budget, error) {
	rows, err := q.db.Query(ctx, getUserBudgets, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Budget
	for rows.Next() {
		var i Budget
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Amount,
			&i.Goal,
			&i.StartDate,
			&i.EndDate,
			&i.UserID,
			&i.CategoryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserBudgetsPaged = `-- name: GetUserBudgetsPaged :many
SELECT id, created_at, updated_at, amount, goal, start_date, end_date, user_id, category_id FROM budgets WHERE user_id = $1
AND created_at >= $2 AND id < $3
ORDER BY created_at ASC, id DESC
LIMIT $4
`

type GetUserBudgetsPagedParams struct {
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Limit     int32     `json:"limit"`
}

func (q *Queries) GetUserBudgetsPaged(ctx context.Context, arg GetUserBudgetsPagedParams) ([]Budget, error) {
	rows, err := q.db.Query(ctx, getUserBudgetsPaged,
		arg.UserID,
		arg.CreatedAt,
		arg.ID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Budget
	for rows.Next() {
		var i Budget
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Amount,
			&i.Goal,
			&i.StartDate,
			&i.EndDate,
			&i.UserID,
			&i.CategoryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBudgetAmount = `-- name: UpdateBudgetAmount :exec
UPDATE budgets SET amount = amount + $2
WHERE category_id = $1 AND start_date <= $3 AND end_date >= $3
`

type UpdateBudgetAmountParams struct {
	CategoryID uuid.UUID       `json:"category_id"`
	Amount     decimal.Decimal `json:"amount"`
	StartDate  time.Time       `json:"start_date"`
}

// Since UpdateBudgetAmount is only called by the API, there is no need to
// check if the user is the owner of the budget since the API already does that
func (q *Queries) UpdateBudgetAmount(ctx context.Context, arg UpdateBudgetAmountParams) error {
	_, err := q.db.Exec(ctx, updateBudgetAmount, arg.CategoryID, arg.Amount, arg.StartDate)
	return err
}
