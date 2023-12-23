-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscribers
(
    id         SERIAL PRIMARY KEY,
    user_id    UUID                     NOT NULL,
    to_user_id UUID                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE          DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE subscribers;
-- +goose StatementEnd
