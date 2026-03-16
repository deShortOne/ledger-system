package dto

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id         int64
	Identifier uuid.UUID
	TransferId int64
	CreatedAt  time.Time
	Status     string
}
