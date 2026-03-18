package domain

import (
	"context"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/google/uuid"
)

type LedgerRepository interface {
	CreateLedgerEntry(ctx context.Context, record dto.LedgerEntry) error
	CreateTransaction(ctx context.Context, record dto.Transaction) error
}

type AccountBalanceRepository interface {
	CreateNewAccountBalance(ctx context.Context, accountId uuid.UUID, updatedAt time.Time) error
	GetAccountBalance(ctx context.Context, accountId uuid.UUID) (dto.AccountBalance, error)
	UpdateAccountBalance(ctx context.Context, record dto.AccountBalance) error
}
