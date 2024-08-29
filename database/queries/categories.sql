-- name: CreateCategory :one
INSERT INTO categories (id, created_at, updated_at, name, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: GetUserCategoriesPaged :many
SELECT * FROM categories WHERE user_id = $1
AND created_at >= $2 AND id < $3
ORDER BY created_at ASC, id DESC
LIMIT $4;

-- name: GetUserCategories :many
SELECT * FROM categories WHERE user_id = $1
ORDER BY created_at ASC, id DESC
LIMIT $2;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;
