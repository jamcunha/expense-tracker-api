-- TODO: add query (and route) to get all expenses in a time interval

-- name: CreateExpense :one
INSERT INTO expenses (id, created_at, updated_at, description, amount, category_id, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteExpense :one
DELETE FROM expenses WHERE id = $1 RETURNING *;

-- name: GetUserExpensesPaged :many
SELECT * FROM expenses WHERE user_id = $1
AND created_at <= $2 AND id < $3
ORDER BY created_at DESC, id DESC
LIMIT $4;

-- name: GetUserExpenses :many
SELECT * FROM expenses WHERE user_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2;

-- name: GetCategoryExpensesPaged :many
SELECT * FROM expenses WHERE category_id = $1
AND created_at <= $2 AND id < $3
ORDER BY created_at DESC, id DESC
LIMIT $4;

-- name: GetCategoryExpenses :many
SELECT * FROM expenses WHERE category_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2;

-- name: GetExpenseByID :one
SELECT * FROM expenses WHERE id = $1;

-- NOTE: both folowing queries are private to the app, not used in the API

-- name: GetTotalSpent :one
SELECT CAST(SUM(amount) AS NUMERIC(10, 4)) FROM expenses
WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3;

-- name: GetTotalSpentInCategory :one
SELECT CAST(SUM(amount) AS NUMERIC(10, 4)) FROM expenses
WHERE user_id = $1 AND category_id = $2 AND created_at >= $3 AND created_at <= $4;
