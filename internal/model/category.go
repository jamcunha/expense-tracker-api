package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/database"
)

// Since the category only interacts with only one user, the model does not need to have a user field.
type Category struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name string `json:"name"`
}

func DatabaseCategoryToCategory(dbCategory database.Category) Category {
	return Category{
		ID:        dbCategory.ID,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
		Name:      dbCategory.Name,
	}
}
