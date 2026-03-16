package dto

import (
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/contracts"
	"github.com/google/uuid"
)

type LedgerEntry struct {
	Id            int64
	Identifier    uuid.UUID
	TransactionId int64
	AccountId     int64
	Amount        float64
	Direction     contracts.LedgerDirection
	CreatedAt     time.Time
}
