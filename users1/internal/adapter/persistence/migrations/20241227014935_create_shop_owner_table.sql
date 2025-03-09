-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_owners (
    user_id UUID PRIMARY KEY,
    bussiness_license VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_owners;
-- +goose StatementEnd
