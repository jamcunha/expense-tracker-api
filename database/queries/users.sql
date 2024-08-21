-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: LoginUser :one
SELECT * FROM users WHERE email = $1 AND password = $2;

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
