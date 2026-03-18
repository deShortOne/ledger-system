package service

import (
	"testing"
	"time"

	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/deshortone/ledger-system/internal/transfer/repository/memory"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatingTransfer(t *testing.T) {
	t.Run("successfully creating transfer", func(t *testing.T) {
		var tx pgx.Tx
		repo := memory.NewTransferInMemoryRepository()
		service := NewTransferService(repo)

		transferRequestId := uuid.New()
		timee, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-15 12:00:00 +0000")
		require.NoError(t, err)
		customTime := dto.NewCustomTime(timee)
		_, err = service.CreateTransfer(t.Context(), tx, transferRequestId, customTime)
		require.NoError(t, err)

		require.Equal(t, 1, len(repo.Transfers))
		assert.Equal(t, transferRequestId, repo.Transfers[0].TransferRequestId)
		assert.Equal(t, customTime, repo.Transfers[0].ExecutedAt)
	})
}

func TestCreatingTransferRequest(t *testing.T) {
	t.Run("successfully creating transfer request", func(t *testing.T) {
		repo := memory.NewTransferInMemoryRepository()
		service := NewTransferService(repo)

		fromAccountId := uuid.New()
		toAccountId := uuid.New()
		timee, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-15 12:00:00 +0000")
		require.NoError(t, err)
		customTime := dto.NewCustomTime(timee)
		_, err = service.CreateTransferRequest(t.Context(), dto.CreateNewTransfer{
			FromAccountId: fromAccountId,
			ToAccountId:   toAccountId,
			Amount:        20,
			RequestedAt:   customTime,
		})
		require.NoError(t, err)

		require.Equal(t, 1, len(repo.TransferRequests))
		assert.Equal(t, fromAccountId, repo.TransferRequests[0].FromAccountId)
		assert.Equal(t, toAccountId, repo.TransferRequests[0].ToAccountId)
		assert.Equal(t, float64(20), repo.TransferRequests[0].Amount)
		assert.Equal(t, customTime, repo.TransferRequests[0].RequestedAt)
	})
}

func TestUpdatingTransferRequestStatus(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		var tx pgx.Tx
		repo := memory.NewTransferInMemoryRepository()
		service := NewTransferService(repo)

		transferRequestId := uuid.New()
		err := service.UpdateTransferRequestStatus(t.Context(), tx, transferRequestId, "status message")
		require.NoError(t, err)

		require.Equal(t, 1, len(repo.StatusUpdates))
		assert.Equal(t, transferRequestId, repo.StatusUpdates[0].Id)
		assert.Equal(t, "status message", repo.StatusUpdates[0].Status)
		assert.Equal(t, false, repo.StatusUpdates[0].WasFailureUpdated)
	})
}

func TestUpdatingTransferRequestStatusWithFailure(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		repo := memory.NewTransferInMemoryRepository()
		service := NewTransferService(repo)

		transferRequestId := uuid.New()
		err := service.UpdateTransferRequestStatusWithFailure(t.Context(), transferRequestId, "status message", "failure message")
		require.NoError(t, err)

		require.Equal(t, 1, len(repo.StatusUpdates))
		assert.Equal(t, transferRequestId, repo.StatusUpdates[0].Id)
		assert.Equal(t, "status message", repo.StatusUpdates[0].Status)
		assert.Equal(t, "failure message", repo.StatusUpdates[0].Failure)
		assert.Equal(t, true, repo.StatusUpdates[0].WasFailureUpdated)
	})
}
