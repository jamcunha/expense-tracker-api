package internal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jamcunha/expense-tracker/internal/repository"
)

// Wrapper for mocking transactions

type DBConn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Close(ctx context.Context) error
}

// Wrapper for queries

type Querier interface {
	repository.Querier
	WithTx(tx pgx.Tx) Querier
}

type Queries struct {
	*repository.Queries
}

func (q *Queries) WithTx(tx pgx.Tx) Querier {
	return &Queries{
		Queries: q.Queries.WithTx(tx),
	}
}

func NewQuerier(q *repository.Queries) Querier {
	return &Queries{
		Queries: q,
	}
}
