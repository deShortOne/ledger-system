package domain

import (
	"context"

	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TransferRepository interface {
	CreateTransfer(ctx context.Context, tx pgx.Tx, request dto.NewTransfer) error
	CreateTransferRequest(ctx context.Context, request dto.NewTransferRequest) error
	UpdateTransferRequestStatusWithTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, status string) error
	UpdateTransferRequestStatusWithFailure(ctx context.Context, id uuid.UUID, status, failure string) error
}
