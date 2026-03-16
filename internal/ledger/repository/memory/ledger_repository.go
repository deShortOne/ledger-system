package memory

import (
	"context"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/jackc/pgx/v5"
)

type LedgerInMemoryRepository struct {
	ledgerEntries []dto.LedgerEntry
	transactions  []dto.Transaction
}

func NewLedgerInMemoryRepository() *LedgerInMemoryRepository {
	return &LedgerInMemoryRepository{
		ledgerEntries: []dto.LedgerEntry{},
		transactions:  []dto.Transaction{},
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
