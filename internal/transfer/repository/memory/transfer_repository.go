package memory

import (
	"context"

	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/google/uuid"
)

type TransferInMemoryRepository struct {
	TransferRequests []dto.NewTransferRequest
	Transfers        []dto.NewTransfer
	StatusUpdates    []TransferRequestStatusUpdate
}

func NewTransferInMemoryRepository() *TransferInMemoryRepository {
	return &TransferInMemoryRepository{}
}

func (r *TransferInMemoryRepository) CreateTransfer(ctx context.Context, request dto.NewTransfer) error {
	r.Transfers = append(r.Transfers, request)
	return nil
}

func (r *TransferInMemoryRepository) CreateTransferRequest(ctx context.Context, request dto.NewTransferRequest) error {
	r.TransferRequests = append(r.TransferRequests, request)
	return nil
}

func (r *TransferInMemoryRepository) UpdateTransferRequestStatusWithTx(ctx context.Context, id uuid.UUID, status string) error {
	r.StatusUpdates = append(r.StatusUpdates, TransferRequestStatusUpdate{
		Id:                id,
		Status:            status,
		WasFailureUpdated: false,
	})
	return nil
}

func (r *TransferInMemoryRepository) UpdateTransferRequestStatusWithFailure(ctx context.Context, id uuid.UUID, status, failure string) error {
	r.StatusUpdates = append(r.StatusUpdates, TransferRequestStatusUpdate{
		Id:                id,
		Status:            status,
		Failure:           failure,
		WasFailureUpdated: true,
	})
	return nil
}

type TransferRequestStatusUpdate struct {
	Id                uuid.UUID
	Status            string
	Failure           string
	WasFailureUpdated bool
}
