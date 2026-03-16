package domain

import (
	"context"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/jackc/pgx/v5"
)

type LedgerRepository interface {
	CreateLedgerEntry(ctx context.Context, tx pgx.Tx, record dto.LedgerEntry) (dto.LedgerEntry, error)
	CreateTransaction(ctx context.Context, tx pgx.Tx, record dto.Transaction) (dto.Transaction, error)
	GetAccountBalance(ctx context.Context, tx pgx.Tx, accountId int64) (dto.AccountBalance, error)
	UpdateAccountBalance(ctx context.Context, tx pgx.Tx, record dto.AccountBalance) error
}
