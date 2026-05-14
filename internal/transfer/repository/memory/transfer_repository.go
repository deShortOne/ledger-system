package memory

import (
	"context"
	"errors"

	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/deshortone/ledger-system/pkg/failure"
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

func (r TransferInMemoryRepository) GetTransferStatus(ctx context.Context, id uuid.UUID) (string, string, error) {
	for _, requestStatusUpdate := range r.StatusUpdates { // should be checking in reverse
		if requestStatusUpdate.Id == id {
			return requestStatusUpdate.Status, requestStatusUpdate.Failure, nil
		}
	}

	return "", "", failure.NewFailure(
		failure.TransferRequestNotFound,
		failure.NotFound,
		errors.New("transfer status update was not found"),
		"No transfer status update record exists for the requested transfer identifier",
	)
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
