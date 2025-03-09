-- +goose Up
-- +goose StatementBegin
CREATE TYPE role_enum AS ENUM ('ADMIN', 'CUSTOMER', 'SHOP_OWNER');
CREATE TABLE users(
    id UUID PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    phone_number VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role role_enum DEFAULT 'CUSTOMER',
    current_refresh_token TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TYPE role_enum;
-- +goose StatementEnd
