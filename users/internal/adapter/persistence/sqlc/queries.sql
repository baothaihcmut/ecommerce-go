-- name: CreateUser :exec
INSERT INTO users(id,email, password ,phone_number, first_name,last_name,role, current_refresh_token)
VALUES (sqlc.narg('id'),
        sqlc.narg('email'),
        sqlc.narg('password'),
        sqlc.narg('phoneNumber'),
        sqlc.narg('firstName'),
        sqlc.narg('lastName'),
        sqlc.narg('role'),
        sqlc.narg('currentRefreshToken')
        );

-- name: UpdateUser :exec
UPDATE users
SET
    email = COALESCE(sqlc.narg('email'),email),
    password = COALESCE(sqlc.narg('password'),password),
    phone_number = COALESCE(sqlc.narg('phoneNumber'),phone_number),
    first_name = COALESCE(sqlc.narg('firstName'),first_name),
    last_name = COALESCE(sqlc.narg('lastName'),last_name),
    role = COALESCE(sqlc.narg('role'),role),
    current_refresh_token = COALESCE(sqlc.narg('currentRefreshToken'),current_refresh_token)
WHERE id = sqlc.narg('id');

-- name: CreateAddress :exec
INSERT INTO addresses(priority,street, town, city, province,user_id)
VALUES (
    sqlc.narg('priority'),
    sqlc.narg('street'),
    sqlc.narg('town'),
    sqlc.narg('city'),
    sqlc.narg('province'),
    sqlc.narg('userId')
);
-- name: UpdateAddress :exec
UPDATE addresses
SET
    street = COALESCE(sqlc.narg('street'),street),
    town = COALESCE(sqlc.narg('town'),town),
    city = COALESCE(sqlc.narg('city'),city),
    province = COALESCE(sqlc.narg('province'),province)
WHERE user_id = sqlc.narg('userId') AND priority = sqlc.narg('priority');


-- name: CreateCustomer :exec
INSERT INTO customers(user_id,loyal_point)
VALUES(
    sqlc.narg('userId'),
    sqlc.narg('loyalPoint')
);

-- name: UpdateCustomer :exec
UPDATE customers
SET
    loyal_point = COALESCE(sqlc.narg('loyalPoint'),loyal_point)
WHERE user_id = sqlc.narg('userId');

-- name: CreateShopOwner :exec
INSERT INTO shop_owners(user_id,bussiness_license)
VALUES(
    sqlc.narg('userId'),
    sqlc.narg('bussinessLicense')
);

-- name: UpdateShopOwner :exec
UPDATE shop_owners
SET
    bussiness_license = COALESCE(sqlc.narg('bussinessLicense'),bussiness_license)
WHERE user_id = sqlc.narg('userId');

-- name: CheckUserExistByCriteria :one
SELECT 1
FROM users u
WHERE
    CASE sqlc.narg('criteria')
        WHEN 'email' THEN u.email = sqlc.narg('value')::text
        WHEN 'phone_number' THEN u.phone_number = sqlc.narg('value')::text
        WHEN 'id' THEN u.id = sqlc.narg('value')::uuid
        WHEN 'firstName' THEN u.first_name = sqlc.narg('value')::text
        WHEN 'lastName' THEN u.last_name = sqlc.narg('value')::text
    END;


-- name: FindAllAddressOfUser :many
SELECT * FROM addresses WHERE user_id = sqlc.narg('userId');
SELECT *
FROM users u
WHERE u.email = sqlc.narg('email')
LIMIT 1;


-- name: FindUserByCriteria :one
SELECT *
FROM users u
LEFT JOIN customers c ON u.id = c.user_id
LEFT JOIN shop_owners s ON u.id = s.user_id
WHERE
    CASE sqlc.narg('criteria')
        WHEN 'email' THEN u.email = sqlc.narg('value')::text
        WHEN 'phone_number' THEN u.phone_number = sqlc.narg('value')::text
        WHEN 'id' THEN u.id = sqlc.narg('value')::uuid
        WHEN 'firstName' THEN u.first_name = sqlc.narg('value')::text
        WHEN 'lastName' THEN u.last_name = sqlc.narg('value')::text
    END
LIMIT 1;

-- name: FindAdminByCriteria :one
SELECT *
FROM admins a
WHERE
    CASE sqlc.narg('criteria')
        WHEN 'email' THEN u.email = sqlc.narg('value')::text
        WHEN 'phone_number' THEN u.phone_number = sqlc.narg('value')::text
        WHEN 'id' THEN u.id = sqlc.narg('value')::uuid
        WHEN 'firstName' THEN u.first_name = sqlc.narg('value')::text
        WHEN 'lastName' THEN u.last_name = sqlc.narg('value')::text
    END
LIMIT 1;

-- name: InsertAdmin :exec
INSERT INTO admins (id, first_name, last_name, email, password, phone_number,current_refresh_token)
				VALUES (
				    sqlc.narg('id'),
					sqlc.narg('firstName'),
					sqlc.narg('lastName'),
					sqlc.narg('email'),
					sqlc.narg('password'),
					sqlc.narg('phoneNumber'),
					sqlc.narg('currentRefreshToken')
				);
-- name: UpdateAdmin :exec
UPDATE admins
SET
    first_name = COALESCE(sqlc.narg('firstName'), first_name),
    last_name = COALESCE(sqlc.narg('lastName'), last_name),
    email = COALESCE(sqlc.narg('email'), email),
    password = COALESCE(sqlc.narg('password'), password),
    phone_number = COALESCE(sqlc.narg('phoneNumber'), phone_number),
    current_refresh_token = COALESCE(sqlc.narg('currentRefreshToken'), current_refresh_token)
WHERE id = sqlc.narg('id');
