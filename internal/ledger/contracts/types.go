package contracts

import "time"

type AddToLedgerRequest struct {
	TransferId int64
	CreatedAt  time.Time
	Entries    []LedgerEntries
}

type LedgerEntries struct {
	AccountId int64
	Amount    float64
	Direction LedgerDirection
}
