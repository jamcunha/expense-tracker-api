package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Budget struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Amount     decimal.Decimal `json:"amount"`
	Goal       decimal.Decimal `json:"goal"`
	StartDate  time.Time       `json:"start_date"`
	EndDate    time.Time       `json:"end_date"`
	UserID     uuid.UUID       `json:"user_id"`
	CategoryID uuid.UUID       `json:"category_id"`
}
