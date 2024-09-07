-- TODO: check on cascade delete on all tables
-- TODO: add a check constraint start_date < end_date

-- +goose Up

CREATE TABLE budgets (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    amount NUMERIC(10, 4) NOT NULL,
    goal NUMERIC(10, 4) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE -- NOTE: here might be better to set NULL instead of deleting

    CONSTRAINT date_check CHECK (start_date < end_date)
);

CREATE INDEX id_budgets_pagination ON expenses (created_at, id);

-- +goose Down

DROP TABLE budgets;
