// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: categories.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO categories (id, created_at, updated_at, name, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, name, user_id
`

type CreateCategoryParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	UserID    uuid.UUID
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	row := q.db.QueryRowContext(ctx, createCategory,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.UserID,
	)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.UserID,
	)
	return i, err
}

const deleteCategory = `-- name: DeleteCategory :one
DELETE FROM categories WHERE id = $1 AND user_id = $2 RETURNING id, created_at, updated_at, name, user_id
`

type DeleteCategoryParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteCategory(ctx context.Context, arg DeleteCategoryParams) (Category, error) {
	row := q.db.QueryRowContext(ctx, deleteCategory, arg.ID, arg.UserID)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.UserID,
	)
	return i, err
}

const getCategoryByID = `-- name: GetCategoryByID :one
SELECT id, created_at, updated_at, name, user_id FROM categories WHERE id = $1 AND user_id = $2
`

type GetCategoryByIDParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) GetCategoryByID(ctx context.Context, arg GetCategoryByIDParams) (Category, error) {
	row := q.db.QueryRowContext(ctx, getCategoryByID, arg.ID, arg.UserID)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.UserID,
	)
	return i, err
}

const getUserCategories = `-- name: GetUserCategories :many
SELECT id, created_at, updated_at, name, user_id FROM categories WHERE user_id = $1
ORDER BY created_at ASC, id DESC
LIMIT $2
`

type GetUserCategoriesParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetUserCategories(ctx context.Context, arg GetUserCategoriesParams) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, getUserCategories, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
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

const getUserCategoriesPaged = `-- name: GetUserCategoriesPaged :many
SELECT id, created_at, updated_at, name, user_id FROM categories WHERE user_id = $1
AND created_at >= $2 AND id < $3
ORDER BY created_at ASC, id DESC
LIMIT $4
`

type GetUserCategoriesPagedParams struct {
	UserID    uuid.UUID
	CreatedAt time.Time
	ID        uuid.UUID
	Limit     int32
}

func (q *Queries) GetUserCategoriesPaged(ctx context.Context, arg GetUserCategoriesPagedParams) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, getUserCategoriesPaged,
		arg.UserID,
		arg.CreatedAt,
		arg.ID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
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
