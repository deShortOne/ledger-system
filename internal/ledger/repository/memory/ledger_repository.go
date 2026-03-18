package memory

import (
	"context"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
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

func (r *LedgerInMemoryRepository) CreateLedgerEntry(ctx context.Context, record dto.LedgerEntry) error {
	r.ledgerEntries = append(r.ledgerEntries, record)
	return nil
}

func (r *LedgerInMemoryRepository) CreateTransaction(ctx context.Context, record dto.Transaction) error {
	r.transactions = append(r.transactions, record)
	return nil
}
