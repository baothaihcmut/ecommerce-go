-- +goose Up
-- +goose StatementBegin
CREATE TYPE role_enum AS ENUM ('admin', 'customer', 'shop_owner');
ALTER TABLE users ADD COLUMN role role_enum DEFAULT 'customer';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE DROP COLUMN role;
-- +goose StatementEnd