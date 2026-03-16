package dto

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Identifier uuid.UUID
	TransferId uuid.UUID
	CreatedAt  time.Time
	Status     string
}
