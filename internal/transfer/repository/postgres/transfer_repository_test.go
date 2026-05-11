package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/deshortone/ledger-system/internal/platform/database_base"
	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	repository := NewTransferPostgresRepository(pool)
	t.Run("success", func(t *testing.T) {
		transferRequestId := uuid.New()
		_, err := pool.Exec(t.Context(), `
		    INSERT INTO transfer.transfer_requests (id, identifier, from_account_id, to_account_id, amount, status, requested_at)
			OVERRIDING SYSTEM VALUE
		    VALUES ($1, $2, $3, $4, $5, $6, NOW());
		`, 1, transferRequestId, 1, 2, 100, "posted")
		require.NoError(t, err)

		transferId := uuid.New()
		executedAtTime, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-15 12:00:00 +0000")
		require.NoError(t, err)

		uow := database_base.NewPgUnitOfWork(pool)
		err = uow.Do(t.Context(), func(ctx context.Context) error {
			return repository.CreateTransfer(ctx, dto.NewTransfer{
				Identifier:        transferId,
				TransferRequestId: transferRequestId,
				ExecutedAt:        dto.NewCustomTime(executedAtTime),
			})
		})
		assert.NoError(t, err)

		status, failureReason, err := repository.GetTransferStatus(t.Context(), transferRequestId)
		require.NoError(t, err)
		assert.Equal(t, "posted", status)
		assert.Equal(t, "", failureReason)

		var executedActual time.Time
		err = pool.QueryRow(t.Context(), "SELECT executed_at FROM transfer.transfers WHERE identifier = $1;", transferId).Scan(&executedActual)
		require.NoError(t, err)
		assert.Equal(t, executedAtTime, executedActual)

		_, err = pool.Exec(t.Context(), `
			DELETE FROM transfer.transfers
			WHERE identifier = $1;
		`, transferId)
		require.NoError(t, err)
		_, err = pool.Exec(t.Context(), `
			DELETE FROM transfer.transfer_requests
			WHERE id = 1;
		`)
		require.NoError(t, err)
	})
}

func TestCreateTransferRequest(t *testing.T) {
	repository := NewTransferPostgresRepository(pool)
	t.Run("success", func(t *testing.T) {
		transferRequestId := uuid.New()
		requestedAtTime, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-17 12:00:00 +0000")
		require.NoError(t, err)

		uow := database_base.NewPgUnitOfWork(pool)
		err = uow.Do(t.Context(), func(ctx context.Context) error {
			return repository.CreateTransferRequest(ctx, dto.NewTransferRequest{
				Identifier:    transferRequestId,
				FromAccountId: account1Id,
				ToAccountId:   account2Id,
				Amount:        100,
				Status:        "la status",
				RequestedAt:   dto.NewCustomTime(requestedAtTime),
			})
		})
		assert.NoError(t, err)

		var identifierActual uuid.UUID
		var fromAccountIdActual int64
		var toAccountIdActual int64
		var amount float64
		var status string
		var requestedAtActual time.Time
		err = pool.QueryRow(t.Context(), `
			SELECT identifier,
				from_account_id,
				to_account_id,
				amount,
				status,
				requested_at
			FROM transfer.transfer_requests
			WHERE identifier = $1;
		`, transferRequestId).Scan(&identifierActual,
			&fromAccountIdActual,
			&toAccountIdActual,
			&amount,
			&status,
			&requestedAtActual)
		require.NoError(t, err)

		assert.Equal(t, transferRequestId, identifierActual)
		assert.Equal(t, int64(1), fromAccountIdActual)
		assert.Equal(t, int64(2), toAccountIdActual)
		assert.Equal(t, float64(100), amount)
		assert.Equal(t, "la status", status)
		assert.Equal(t, requestedAtTime, requestedAtActual)

		_, err = pool.Exec(t.Context(), `
			DELETE FROM transfer.transfer_requests
			WHERE identifier = $1;
		`, transferRequestId)
		require.NoError(t, err)
	})
}
