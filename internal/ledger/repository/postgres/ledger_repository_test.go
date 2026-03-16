package postgres

import (
	"testing"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreatingLedgerEntry(t *testing.T) {
	repository := NewLedgerPostgresRepository(pool)

	t.Run("when transaction is not aborted", func(t *testing.T) {
		tx, err := pool.Begin(t.Context())
		if err != nil {
			panic(err)
		}
		defer tx.Rollback(t.Context())

		transactionId := uuid.New()
		err = repository.CreateTransaction(t.Context(), tx, dto.Transaction{
			Identifier: transactionId,
			TransferId: transferId,
			CreatedAt:  time.Now(),
			Status:     "pending",
		})
		require.NoError(t, err)

		entryId := uuid.New()
		err = repository.CreateLedgerEntry(t.Context(), tx, dto.LedgerEntry{
			Identifier:    entryId,
			TransactionId: transactionId,
			AccountId:     account1Id,
			Amount:        100,
			Direction:     "CREDIT",
			CreatedAt:     time.Now(),
		})
		require.NoError(t, err)

		require.NoError(t, tx.Commit(t.Context()))

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
		tx, err := pool.Begin(t.Context())
		if err != nil {
			panic(err)
		}

		transactionId := uuid.New()
		err = repository.CreateTransaction(t.Context(), tx, dto.Transaction{
			Identifier: transactionId,
			TransferId: transferId,
			CreatedAt:  time.Now(),
			Status:     "pending",
		})
		require.NoError(t, err)

		entryId := uuid.New()
		err = repository.CreateLedgerEntry(t.Context(), tx, dto.LedgerEntry{
			Identifier:    entryId,
			TransactionId: transactionId,
			AccountId:     account1Id,
			Amount:        100,
			Direction:     "CREDIT",
			CreatedAt:     time.Now(),
		})
		require.NoError(t, err)

		require.NoError(t, tx.Rollback(t.Context()))

		// verify row exists, etc.
		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.ledger_entries WHERE identifier = $1", entryId).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}

func TestCreatingTransaction(t *testing.T) {
	repository := NewLedgerPostgresRepository(pool)

	t.Run("when transaction is not aborted", func(t *testing.T) {
		tx, err := pool.Begin(t.Context())
		if err != nil {
			panic(err)
		}
		defer tx.Rollback(t.Context())

		transactionId := uuid.New()
		err = repository.CreateTransaction(t.Context(), tx, dto.Transaction{
			Identifier: transactionId,
			TransferId: transferId,
			CreatedAt:  time.Now(),
			Status:     "pending",
		})
		require.NoError(t, err)

		err = tx.Commit(t.Context())
		require.NoError(t, err)

		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.transactions WHERE identifier = $1", transactionId).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)

		// clean up
		_, err = pool.Exec(t.Context(), "TRUNCATE TABLE ledger.transactions;")
	})

	t.Run("when transaction is aborted", func(t *testing.T) {
		tx, err := pool.Begin(t.Context())
		if err != nil {
			panic(err)
		}

		transactionId := uuid.New()
		err = repository.CreateTransaction(t.Context(), tx, dto.Transaction{
			Identifier: transactionId,
			TransferId: transferId,
			CreatedAt:  time.Now(),
			Status:     "pending",
		})
		require.NoError(t, err)

		err = tx.Rollback(t.Context())
		require.NoError(t, err)

		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.transactions WHERE identifier = $1", transactionId).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}
