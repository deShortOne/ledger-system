-- name: CreateAccount :exec
INSERT INTO identity.accounts (identifier, user_id, created_at, account_type, currency, status)
SELECT  $1, users.id, $3, $4, $5, $6
FROM identity.users
WHERE users.identifier = $2;

-- name: GetAccountsOwnedByUser :many
SELECT id,
    identifier account_identifier,
    created_at,
    account_type,
    currency,
    status
FROM identity.accounts
WHERE user_id = (
    SELECT id
    FROM identity.users
    WHERE users.identifier = $1
);
