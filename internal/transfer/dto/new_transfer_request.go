package dto

import (
	"github.com/google/uuid"
)

type NewTransferRequest struct {
	Identifier    uuid.UUID
	FromAccountId uuid.UUID
	ToAccountId   uuid.UUID
	Amount        float64
	Status        string
	RequestedAt   CustomTime
}
