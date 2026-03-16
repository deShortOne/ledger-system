package contracts

import (
	"time"

	"github.com/google/uuid"
)

type AddToLedgerRequest struct {
	TransferId uuid.UUID
	CreatedAt  time.Time
	Entries    []LedgerEntries
}

type LedgerEntries struct {
	AccountId uuid.UUID
	Amount    float64
	Direction LedgerDirection
}
