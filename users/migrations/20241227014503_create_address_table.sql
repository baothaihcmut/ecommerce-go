-- +goose Up
-- +goose StatementBegin
CREATE TABLE addresses (
    priority INT NOT NULL,
    street VARCHAR(100) NOT NULL,
    town VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (priority, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE addresses;
-- +goose StatementEnd
