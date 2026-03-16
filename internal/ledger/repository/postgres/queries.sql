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
