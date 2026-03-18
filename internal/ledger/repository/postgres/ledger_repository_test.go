package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/deshortone/ledger-system/internal/platform/database_base"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreatingLedgerEntry(t *testing.T) {
	repository := NewLedgerPostgresRepository(pool)

	t.Run("when transaction is not aborted", func(t *testing.T) {
		var err error

		uow := database_base.NewPgUnitOfWork(pool)
		transactionId := uuid.New()
		err = uow.Do(t.Context(), func(ctx context.Context) error {
			return repository.CreateTransaction(ctx, dto.Transaction{
				Identifier: transactionId,
				TransferId: transferId,
				CreatedAt:  time.Now(),
				Status:     "pending",
			})
		})
		require.NoError(t, err)

		entryId := uuid.New()
		err = uow.Do(t.Context(), func(ctx context.Context) error {
			return repository.CreateLedgerEntry(ctx, dto.LedgerEntry{
				Identifier:    entryId,
				TransactionId: transactionId,
				AccountId:     account1Id,
				Amount:        100,
				Direction:     "CREDIT",
				CreatedAt:     time.Now(),
			})
		})
		require.NoError(t, err)

		// verify row exists, etc.
		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.ledger_entries WHERE identifier = $1", entryId).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)

		// clean up
		_, err = pool.Exec(t.Context(), "TRUNCATE TABLE ledger.ledger_entries;")
		require.NoError(t, err)
		_, err = pool.Exec(t.Context(), "DELETE FROM ledger.transactions where 1011 = 1011;")
		require.NoError(t, err)
	})

	t.Run("when transaction is aborted", func(t *testing.T) {
		var err error

		transactionId := uuid.New()
		entryId := uuid.New()
		uow := database_base.NewPgUnitOfWork(pool)
		err = uow.Do(t.Context(), func(ctx context.Context) error {
			err = repository.CreateTransaction(ctx, dto.Transaction{
				Identifier: transactionId,
				TransferId: transferId,
				CreatedAt:  time.Now(),
				Status:     "pending",
			})
			require.NoError(t, err)

			err = repository.CreateLedgerEntry(ctx, dto.LedgerEntry{
				Identifier:    entryId,
				TransactionId: transactionId,
				AccountId:     account1Id,
				Amount:        100,
				Direction:     "CREDIT",
				CreatedAt:     time.Now(),
			})
			require.NoError(t, err)

			return errors.New("something to ensure uow rolls this back")
		})
		require.Error(t, err)

		// verify row exists, etc.
		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.ledger_entries WHERE identifier = $1", entryId).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}

func TestCreatingTransaction(t *testing.T) {
	var err error
	repository := NewLedgerPostgresRepository(pool)

	t.Run("when transaction is not aborted", func(t *testing.T) {
		transactionId := uuid.New()

		uow := database_base.NewPgUnitOfWork(pool)
		err = uow.Do(t.Context(), func(ctx context.Context) error {
			return repository.CreateTransaction(ctx, dto.Transaction{
				Identifier: transactionId,
				TransferId: transferId,
				CreatedAt:  time.Now(),
				Status:     "pending",
			})
		})
		require.NoError(t, err)

		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.transactions WHERE identifier = $1", transactionId).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)

		// clean up
		_, err = pool.Exec(t.Context(), "TRUNCATE TABLE ledger.transactions;")
	})

	t.Run("when transaction is aborted", func(t *testing.T) {
		transactionId := uuid.New()
		uow := database_base.NewPgUnitOfWork(pool)
		err = uow.Do(t.Context(), func(ctx context.Context) error {
			err = repository.CreateTransaction(ctx, dto.Transaction{
				Identifier: transactionId,
				TransferId: transferId,
				CreatedAt:  time.Now(),
				Status:     "pending",
			})
			require.NoError(t, err)

			return errors.New("throw fake")
		})
		require.Error(t, err)

		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.transactions WHERE identifier = $1", transactionId).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}
