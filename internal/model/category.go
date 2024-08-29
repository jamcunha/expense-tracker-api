package model

import (
	"time"

	"github.com/google/uuid"
)

// Since the category only interacts with only one user, the model does not need to have a user field.
type Category struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name   string    `json:"name"`
	UserID uuid.UUID `json:"user_id"`
}
