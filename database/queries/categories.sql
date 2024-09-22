-- name: CreateCategory :one
INSERT INTO categories (id, created_at, updated_at, name, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteCategory :one
DELETE FROM categories WHERE id = $1 AND user_id = $2 RETURNING *;

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
SELECT * FROM categories WHERE id = $1 AND user_id = $2;

-- name: UpdateCategory :one
UPDATE categories SET name = $1 AND updated_at = $2
WHERE id = $3 AND user_id = $4 RETURNING *;
