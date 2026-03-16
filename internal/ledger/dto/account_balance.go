package dto

import (
	"time"

	"github.com/google/uuid"
)

type AccountBalance struct {
	AccountId        uuid.UUID
	Availablebalance float64
	UpdatedAt        time.Time
}
