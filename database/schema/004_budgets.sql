-- TODO: check on cascade delete on all tables

-- +goose Up

CREATE TABLE budgets (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    amount NUMERIC(10, 4) NOT NULL,
    goal NUMERIC(10, 4) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    category_id UUID NOT NULL REFERENCES categories(id)
);

CREATE INDEX id_budgets_pagination ON expenses (created_at, id);

-- +goose Down

DROP TABLE budgets;
