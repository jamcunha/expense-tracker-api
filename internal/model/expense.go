package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Expense struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	CategoryID  uuid.UUID       `json:"category_id"`
	UserID      uuid.UUID       `json:"user_id"`
}
