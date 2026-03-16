package dto

import (
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/contracts"
	"github.com/google/uuid"
)

type LedgerEntry struct {
	Identifier    uuid.UUID
	TransactionId uuid.UUID
	AccountId     uuid.UUID
	Amount        float64
	Direction     contracts.LedgerDirection
	CreatedAt     time.Time
}
