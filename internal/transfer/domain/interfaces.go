package domain

import (
	"context"

	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/google/uuid"
)

type TransferRepository interface {
	CreateTransfer(ctx context.Context, request dto.NewTransfer) error
	CreateTransferRequest(ctx context.Context, request dto.NewTransferRequest) error
	UpdateTransferRequestStatusWithTx(ctx context.Context, id uuid.UUID, status string) error
	UpdateTransferRequestStatusWithFailure(ctx context.Context, id uuid.UUID, status, failure string) error
}

type TransferService interface {
	CreateTransferRequest(ctx context.Context, request dto.CreateNewTransfer) (uuid.UUID, error)
	CreateTransfer(ctx context.Context, transferRequestId uuid.UUID, executedAt dto.CustomTime) (uuid.UUID, error)
	UpdateTransferRequestStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateTransferRequestStatusWithFailure(ctx context.Context, id uuid.UUID, status, failure string) error
}
