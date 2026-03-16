-- name: CreateLedgerEntry :exec
INSERT INTO ledger.ledger_entries (identifier, transaction_id, account_id, amount, direction, created_at)
VALUES (
    $1,
    (SELECT transactions.id FROM ledger.transactions where transactions.identifier = $2),
    (SELECT accounts.id FROM identity.accounts where accounts.identifier = $3),
    $4,
    $5,
    $6
);

-- name: CreateTransaction :exec
INSERT INTO ledger.transactions (identifier, transfer_id, created_at, status) 
VALUES (
    $1,
    (SELECT transfers.id FROM transfer.transfers WHERE transfers.identifier = $2),
    $3,
    $4
);

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
