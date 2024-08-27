-- name: CreateCategory :one
INSERT INTO categories (id, created_at, updated_at, name, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1 AND user_id = $2;

-- name: GetUserCategories :many
SELECT * FROM categories WHERE user_id = $1;
