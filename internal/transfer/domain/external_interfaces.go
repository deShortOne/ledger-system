package domain

import (
	"context"

	"github.com/deshortone/ledger-system/internal/ledger/contracts"
	"github.com/jackc/pgx/v5"
)

type LedgerService interface {
	AddToLedger(ctx context.Context, tx pgx.Tx, request contracts.AddToLedgerRequest) error
}
