package dto

import (
	"github.com/google/uuid"
)

type CreateNewTransfer struct {
	FromAccountId uuid.UUID
	ToAccountId   uuid.UUID
	Amount        float64
	RequestedAt   CustomTime
}
