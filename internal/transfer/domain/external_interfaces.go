package domain

import (
	"context"

	"github.com/deshortone/ledger-system/internal/ledger/contracts"
)

type LedgerService interface {
	AddToLedger(ctx context.Context, request contracts.AddToLedgerRequest) error
}
