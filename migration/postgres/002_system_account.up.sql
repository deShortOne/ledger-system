INSERT INTO identity.users (identifier, first_name, last_name) VALUES
('31a105a5-855e-45a3-856b-d22013a2c69d', 'SYSTEM', 'SYSTEM');

INSERT INTO identity.accounts (identifier, user_id, created_at, account_type, currency, status) VALUES
('6724081d-6f50-4172-92c1-9c5d571f051c', 1, '2026-05-10T15:40:00Z', 'WORLD', 'GBP', 'available');


INSERT INTO transfer.transfer_requests(
    identifier, from_account_id, to_account_id, amount, status, failure_reason, requested_at) VALUES
('fc34a724-f0a7-4782-a814-9496182078dc', 1, 1, 100000000000, 'Success', null, '2026-05-10T15:40:00Z'); -- yes, transferring to the same account

INSERT INTO transfer.transfers(identifier, transfer_request_id, executed_at) VALUES
('a056caec-e826-4ccd-bb10-d173d45bb88b', 1, '2026-05-10T15:40:00Z');


INSERT INTO ledger.transactions(identifier, transfer_id, created_at, status) VALUES
('18cebc4f-7221-45be-b818-a061a5094b7b', 1, '2026-05-10T15:40:00Z', 'posted');

INSERT INTO ledger.ledger_entries(
	identifier, transaction_id, account_id, amount, direction, created_at) VALUES
('fb29c34c-53a9-47f1-9b7b-06dda42845df', 1, 1, 100000000000, 'CREDIT', '2026-05-10T15:40:00Z');

INSERT INTO ledger.account_balances(account_id, available_balance, updated_at) VALUES
(1, 100000000000, '2026-05-10T15:40:00Z');
