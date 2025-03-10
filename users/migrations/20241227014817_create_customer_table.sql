-- +goose Up
-- +goose StatementBegin
CREATE TABLE customers (
    user_id UUID PRIMARY KEY,
    loyal_point INT DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE customers;
-- +goose StatementEnd
