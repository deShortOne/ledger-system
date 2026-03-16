-- name: CreateLedgerEntry :one
INSERT INTO ledger.ledger_entries (
    identifier, transaction_id, account_id, amount, direction, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id;

-- name: CreateTransaction :one
INSERT INTO ledger.transactions (
    identifier, transfer_id, created_at, status
) VALUES (
    $1, $2, $3, $4
)
RETURNING id;

-- name: GetAccountBalanceAndLock :one
SELECT available_balance, updated_at
FROM ledger.account_balances
WHERE account_id = $1
FOR UPDATE;

-- name: UpdateAccountBalance :exec
UPDATE ledger.account_balances
SET available_balance = $2,
    updated_at = $3
WHERE account_id = $1;
