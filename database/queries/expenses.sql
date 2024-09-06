-- TODO: Order expenses by created_at DESC (to test with cursor pagination)

-- name: CreateExpense :one
INSERT INTO expenses (id, created_at, updated_at, description, amount, category_id, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses WHERE id = $1;

-- name: GetUserExpensesPaged :many
SELECT * FROM expenses WHERE user_id = $1
AND created_at >= $2 AND id < $3
ORDER BY created_at ASC, id DESC
LIMIT $4;

-- name: GetUserExpenses :many
SELECT * FROM expenses WHERE user_id = $1
ORDER BY created_at ASC, id DESC
LIMIT $2;

-- name: GetCategoryExpensesPaged :many
SELECT * FROM expenses WHERE category_id = $1
AND created_at >= $2 AND id < $3
ORDER BY created_at ASC, id DESC
LIMIT $4;

-- name: GetCategoryExpenses :many
SELECT * FROM expenses WHERE category_id = $1
ORDER BY created_at ASC, id DESC
LIMIT $2;

-- name: GetExpenseById :one
SELECT * FROM expenses WHERE id = $1;