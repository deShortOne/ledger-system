package memory

import (
	"context"
	"errors"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LedgerInMemoryRepository struct {
	ledgerEntries   []dto.LedgerEntry
	transactions    []dto.Transaction
	accountBalances []dto.AccountBalance
}

func NewLedgerInMemoryRepository() *LedgerInMemoryRepository {
	return &LedgerInMemoryRepository{
		ledgerEntries:   []dto.LedgerEntry{},
		transactions:    []dto.Transaction{},
		accountBalances: []dto.AccountBalance{},
	}
}

func (r *LedgerInMemoryRepository) CreateLedgerEntry(ctx context.Context, tx pgx.Tx, record dto.LedgerEntry) error {
	r.ledgerEntries = append(r.ledgerEntries, record)
	return nil
}

func (r *LedgerInMemoryRepository) CreateTransaction(ctx context.Context, tx pgx.Tx, record dto.Transaction) error {
	r.transactions = append(r.transactions, record)
	return nil
}

func (r *LedgerInMemoryRepository) GetAccountBalance(ctx context.Context, tx pgx.Tx, accountId uuid.UUID) (dto.AccountBalance, error) {
	for _, accountBalance := range r.accountBalances {
		if accountBalance.AccountId == accountId {
			return accountBalance, nil
		}
	}

	return dto.AccountBalance{}, errors.New("account balance not found")
}

func (r *LedgerInMemoryRepository) UpdateAccountBalance(ctx context.Context, tx pgx.Tx, record dto.AccountBalance) error {
	for i, accountBalance := range r.accountBalances {
		if accountBalance.AccountId == record.AccountId {
			r.accountBalances[i].Availablebalance = record.Availablebalance
			r.accountBalances[i].UpdatedAt = record.UpdatedAt
			return nil
		}
	}

	return errors.New("account balance not found")
}

func (r *LedgerInMemoryRepository) CreateAccountBalance(ctx context.Context, tx pgx.Tx, record dto.AccountBalance) error {
	r.accountBalances = append(r.accountBalances, record)
	return nil
}
