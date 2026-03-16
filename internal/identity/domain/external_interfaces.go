package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccountBalanceCreator interface {
	CreateNewAccount(ctx context.Context, accountId uuid.UUID, createdAt time.Time) error
}
