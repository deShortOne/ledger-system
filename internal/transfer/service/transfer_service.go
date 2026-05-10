package service

import (
	"context"

	"github.com/deshortone/ledger-system/internal/transfer/domain"
	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/google/uuid"
)

type TransferService struct {
	repository domain.TransferRepository
}

func NewTransferService(repository domain.TransferRepository) *TransferService {
	return &TransferService{
		repository: repository,
	}
}
func (t *TransferService) CreateTransfer(ctx context.Context, transferRequestId uuid.UUID, executedAt dto.CustomTime) (uuid.UUID, error) {
	transferId := uuid.New()
	err := t.repository.CreateTransfer(ctx, dto.NewTransfer{
		Identifier:        transferId,
		TransferRequestId: transferRequestId,
		ExecutedAt:        executedAt,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return transferId, nil
}

func (t *TransferService) CreateTransferRequest(ctx context.Context, request dto.CreateNewTransfer) (uuid.UUID, error) {
	transferRequestId := uuid.New()
	err := t.repository.CreateTransferRequest(ctx, dto.NewTransferRequest{
		Identifier:    transferRequestId,
		FromAccountId: request.FromAccountId,
		ToAccountId:   request.ToAccountId,
		Amount:        request.Amount,
		Status:        "pending",
		RequestedAt:   request.RequestedAt,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return transferRequestId, err
}

func (t *TransferService) UpdateTransferRequestStatus(ctx context.Context, id uuid.UUID, status string) error {
	return t.repository.UpdateTransferRequestStatusWithTx(ctx, id, status)
}

func (t *TransferService) UpdateTransferRequestStatusWithFailure(ctx context.Context, id uuid.UUID, status, failure string) error {
	return t.repository.UpdateTransferRequestStatusWithFailure(ctx, id, status, failure)
}
