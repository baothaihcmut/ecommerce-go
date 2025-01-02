-- name: CreateUser :exec
INSERT INTO users(id,email, phone_number, first_name,last_name)
VALUES (sqlc.narg('id'),
        sqlc.narg('email'),
        sqlc.narg('phoneNumber'),
        sqlc.narg('firstName'),
        sqlc.narg('lastName')
        );

-- name: CreateAddress :exec
INSERT INTO addresses(street, town, city, province,user_id)
VALUES (
    sqlc.narg('street'),
    sqlc.narg('town'),
    sqlc.narg('city'),
    sqlc.narg('province'),
    sqlc.narg('userId')
);

-- name: CreateCustomer :exec
INSERT INTO customers(user_id,loyal_point)
VALUES(
    sqlc.narg('userId'),
    sqlc.narg('loyalPoint')
);

-- name: CreateShopOwner :exec
INSERT INTO shop_owners(user_id,bussiness_license)
VALUES(
    sqlc.narg('userId'),
    sqlc.narg('bussinessLicense')
);

-- name: FindUserByCriteria :one
SELECT *
FROM users u
LEFT JOIN customers c ON u.id = c.user_id
LEFT JOIN shop_owners s ON u.id = s.user_id
WHERE
    (
        (sqlc.narg('criteria') = 'email' AND u.email = sqlc.narg('value')) OR
        (sqlc.narg('criteria') = 'phoneNumber' AND u.phone_number = sqlc.narg('value')) OR
        (sqlc.narg('criteria') = 'id' AND u.id = sqlc.narg('value')) OR
        (sqlc.narg('criteria') = 'firstName' AND u.first_name = sqlc.narg('value')) OR
        (sqlc.narg('criteria') = 'lastName' AND u.last_name = sqlc.narg('value'))
    )
LIMIT 1;


-- name: FindAllAddressOfUser :many
SELECT * FROM addresses WHERE user_id = sqlc.narg('userId');
