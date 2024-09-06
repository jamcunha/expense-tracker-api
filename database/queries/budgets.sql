-- name: CreateBudget :one
INSERT INTO budgets (id, created_at, updated_at, amount, goal, start_date, end_date, user_id, category_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateBudgetAmount :exec
UPDATE budgets SET amount = amount + $2
WHERE category_id = $1 AND start_date <= $3 AND end_date >= $3;

-- name: DeleteBudget :exec
DELETE FROM budgets WHERE id = $1;

-- name: GetUserBudgetsPaged :many
SELECT * FROM budgets WHERE user_id = $1
AND created_at >= $2 AND id < $3
ORDER BY created_at ASC, id DESC
LIMIT $4;

-- name: GetUserBudgets :many
SELECT * FROM budgets WHERE user_id = $1
ORDER BY created_at ASC, id DESC
LIMIT $2;

-- name: GetBudgetByID :one
SELECT * FROM budgets WHERE id = $1;
