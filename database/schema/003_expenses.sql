-- +goose Up

CREATE TABLE expenses (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    description TEXT NOT NULL,
    amount NUMERIC(10, 4) NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE, -- NOTE: here might be better to set NULL instead of deleting
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX id_expenses_pagination ON expenses (created_at, id);

-- +goose Down

DROP TABLE expenses;
