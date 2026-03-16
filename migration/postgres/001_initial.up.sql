CREATE TYPE ledger_entry_direction AS ENUM ('CREDIT', 'DEBIT');

CREATE SCHEMA identity;

CREATE TABLE identity.users (
	id BIGINT GENERATED ALWAYS AS IDENTITY,
	identifier UUID NOT NULL,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	PRIMARY KEY(id),
	UNIQUE(identifier)
);

CREATE TABLE identity.accounts (
	id BIGINT GENERATED ALWAYS AS IDENTITY,
	identifier UUID NOT NULL,
	user_id BIGINT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	account_type TEXT NOT NULL,
	currency CHAR(3) NOT NULL,
	status TEXT NOT NULL,
	PRIMARY KEY(id),
	UNIQUE(identifier),
	CONSTRAINT fk_accounts_to_user
		FOREIGN KEY(user_id)
			REFERENCES identity.users(id)
);

CREATE SCHEMA transfer;

CREATE TABLE transfer.transfer_requests (
	id BIGINT GENERATED ALWAYS AS IDENTITY,
	identifier UUID NOT NULL,
	from_account_id  BIGINT NOT NULL,
	to_account_id BIGINT NOT NULL,
	amount NUMERIC(18, 6) NOT NULL,
	status TEXT NOT NULL,
	failure_reason TEXT NULL,
	requested_at TIMESTAMP WITH TIME ZONE NOT NULL,
	PRIMARY KEY(id),
	UNIQUE(identifier),
	CONSTRAINT fk_transfer_requests_from_account_to_account
		FOREIGN KEY(from_account_id)
			REFERENCES identity.accounts(id),
	CONSTRAINT fk_transfer_requests_to_account_to_account
		FOREIGN KEY(to_account_id)
			REFERENCES identity.accounts(id)
);

CREATE TABLE transfer.transfers (
	id BIGINT GENERATED ALWAYS AS IDENTITY,
	identifier UUID NOT NULL,
	transfer_request_id BIGINT NOT NULL,
	executed_at TIMESTAMP WITH TIME ZONE NOT NULL,
	PRIMARY KEY(id),
	UNIQUE (identifier),
	CONSTRAINT fk_transfers_to_transfer_requests
		FOREIGN KEY(transfer_request_id)
			REFERENCES transfer.transfer_requests(id)
);

CREATE SCHEMA ledger;

CREATE TABLE ledger.account_balances (
	account_id BIGINT NOT NULL,
	available_balance NUMERIC(18, 6) NOT NULL,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
	CONSTRAINT fk_account_balances_to_account
		FOREIGN KEY(account_id)
			REFERENCES identity.accounts(id)
);

CREATE TABLE ledger.transactions (
	id BIGINT GENERATED ALWAYS AS IDENTITY,
	identifier UUID NOT NULL,
	transfer_id BIGINT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	status TEXT NOT NULL,
	PRIMARY KEY(id),
	UNIQUE(identifier),
	CONSTRAINT fk_transactions_transfer_id_to_transfer
		FOREIGN KEY(transfer_id)
			REFERENCES transfer.transfers(id)
);

CREATE TABLE ledger.ledger_entries (
	id BIGINT GENERATED ALWAYS AS IDENTITY,
	identifier UUID NOT NULL,
	transaction_id BIGINT NOT NULL,
	account_id BIGINT NOT NULL,
	amount NUMERIC(18, 6) NOT NULL,
	direction ledger_entry_direction NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	PRIMARY KEY(id),
	UNIQUE(identifier),
	CONSTRAINT fk_ledger_entry_to_transaction
		FOREIGN KEY(transaction_id)
			REFERENCES ledger.transactions(id),
	CONSTRAINT fk_ledger_entry_to_account
		FOREIGN KEY(account_id)
			REFERENCES identity.accounts(id)
);
