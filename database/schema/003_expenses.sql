-- +goose Up

CREATE TABLE expenses (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    description TEXT NOT NULL,
    amount NUMERIC(10, 4) NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id),
    user_id UUID NOT NULL REFERENCES users(id)
);

CREATE INDEX id_expenses_pagination ON expenses (created_at, id);

-- +goose Down

DROP TABLE expenses;
