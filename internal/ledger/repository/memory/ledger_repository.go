package memory

import (
	"context"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
)

type LedgerInMemoryRepository struct {
	LedgerEntries []dto.LedgerEntry
	Transactions  []dto.Transaction
}

func NewLedgerInMemoryRepository() *LedgerInMemoryRepository {
	return &LedgerInMemoryRepository{
		LedgerEntries: []dto.LedgerEntry{},
		Transactions:  []dto.Transaction{},
	}
}

func (r *LedgerInMemoryRepository) CreateLedgerEntry(ctx context.Context, record dto.LedgerEntry) error {
	r.LedgerEntries = append(r.LedgerEntries, record)
	return nil
}

func (r *LedgerInMemoryRepository) CreateTransaction(ctx context.Context, record dto.Transaction) error {
	r.Transactions = append(r.Transactions, record)
	return nil
}
