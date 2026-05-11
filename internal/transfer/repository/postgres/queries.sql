-- name: CreateTransferRequest :exec
INSERT INTO transfer.transfer_requests (
	identifier, from_account_id, to_account_id, amount, status, requested_at
) VALUES (
    $1,
    (SELECT accounts.id FROM identity.accounts WHERE accounts.identifier = $2),
    (SELECT accounts.id FROM identity.accounts WHERE accounts.identifier = $3),
    $4,
    $5,
    $6
);

-- name: UpdateTransferRequestStatus :exec
UPDATE transfer.transfer_requests
SET status = $2,
    failure_reason = $3
WHERE identifier = $1;

-- name: CreateTransfer :exec
INSERT INTO transfer.transfers (
	identifier, transfer_request_id, executed_at
) VALUES (
    $1,
    (SELECT transfer_requests.id FROM transfer.transfer_requests WHERE transfer_requests.identifier = $2),
    $3
);

-- name: GetTranserRequestStatus :one
SELECT status, failure_reason
FROM transfer.transfer_requests
WHERE identifier = $1;
