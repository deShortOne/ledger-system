-- name: GetAccountBalanceAndLock :one
SELECT available_balance, updated_at
FROM ledger.account_balances
WHERE account_id = (
    SELECT accounts.id
    FROM identity.accounts
    WHERE accounts.identifier = $1
)
FOR UPDATE;

-- name: UpdateAccountBalance :exec
UPDATE ledger.account_balances
SET available_balance = $2,
    updated_at = $3
WHERE account_id = (
    SELECT accounts.id
    FROM identity.accounts
    WHERE accounts.identifier = $1
);

-- name: CreateAccountBalance :exec
INSERT INTO ledger.account_balances (
    account_id, available_balance, updated_at
) VALUES (
    (SELECT accounts.id FROM identity.accounts WHERE accounts.identifier = $1),
    $2,
    $3
);
