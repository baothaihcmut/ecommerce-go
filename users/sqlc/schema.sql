CREATE TABLE users(
    id UUID PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    is_shop_owner_active BOOLEAN DEFAULT FALSE NOT NULL,
    phone_number VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    current_refresh_token TEXT
);

CREATE TABLE addresses (
    priority INT NOT NULL,
    street VARCHAR(100) NOT NULL,
    town VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (priority, user_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE customers (
    user_id UUID PRIMARY KEY,
    loyal_point INT DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE shop_owners (
    user_id UUID PRIMARY KEY,
    bussiness_license VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE admins (
    id UUID PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    phone_number VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    current_refresh_token TEXT
);
