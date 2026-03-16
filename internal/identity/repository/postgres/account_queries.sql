-- name: CreateAccount :one
INSERT INTO identity.accounts (
    identifier,
    user_id,
    created_at,
    account_type,
    currency,
    "status"
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id;

-- name: GetAccountsOwnedByUser :many
SELECT id, identifier, created_at, account_type, currency, status
FROM identity.accounts
WHERE user_id = $1;
