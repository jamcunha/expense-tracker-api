-- name: CreateBudget :one
INSERT INTO budgets (id, created_at, updated_at, amount, goal, start_date, end_date, user_id, category_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: DeleteBudget :one
DELETE FROM budgets WHERE id = $1 AND user_id = $2 RETURNING *;

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
SELECT * FROM budgets WHERE id = $1 AND user_id = $2;

-- name: UpdateBudgetAmount :exec
-- Since UpdateBudgetAmount is only called by the API, there is no need to
-- check if the user is the owner of the budget since the API already does that
UPDATE budgets SET amount = amount + $2
WHERE category_id = $1 AND start_date <= $3 AND end_date >= $3;
