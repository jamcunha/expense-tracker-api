-- +goose Up

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id)
);

-- +goose Down

DROP TABLE categories;
