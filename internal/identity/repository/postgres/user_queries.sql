-- name: CreateUser :exec
INSERT INTO identity.users (
    identifier,
    first_name,
    last_name
) VALUES (
    $1, $2, $3
);

-- name: GetUser :one
SELECT id, first_name, last_name
FROM identity.users
WHERE identifier = $1;
