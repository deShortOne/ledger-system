package dto

import "github.com/google/uuid"

type NewTransfer struct {
	Identifier        uuid.UUID
	TransferRequestId uuid.UUID
	ExecutedAt        CustomTime
}
