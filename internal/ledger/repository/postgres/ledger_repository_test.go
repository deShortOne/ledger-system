package postgres

import (
	"testing"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

		txn, err := repository.CreateTransaction(t.Context(), tx, dto.Transaction{
			Identifier: uuid.New(),
			TransferId: 1,
			CreatedAt:  time.Now(),
			Status:     "pending",
		})
		require.NoError(t, err)

		entry, err := repository.CreateLedgerEntry(t.Context(), tx, dto.LedgerEntry{
			Identifier:    uuid.New(),
			TransactionId: txn.Id,
			AccountId:     1,
			Amount:        100,
			Direction:     "CREDIT",
			CreatedAt:     time.Now(),
		})
		require.NoError(t, err)

		require.NoError(t, tx.Commit(t.Context()))

		// verify row exists, etc.
		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.ledger_entries WHERE id=$1", entry.Id).Scan(&count)
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

		txn, err := repository.CreateTransaction(t.Context(), tx, dto.Transaction{
			Identifier: uuid.New(),
			TransferId: 1,
			CreatedAt:  time.Now(),
			Status:     "pending",
		})
		require.NoError(t, err)

		entry, err := repository.CreateLedgerEntry(t.Context(), tx, dto.LedgerEntry{
			Identifier:    uuid.New(),
			TransactionId: txn.Id,
			AccountId:     1,
			Amount:        100,
			Direction:     "CREDIT",
			CreatedAt:     time.Now(),
		})
		require.NoError(t, err)

		require.NoError(t, tx.Rollback(t.Context()))

		// verify row exists, etc.
		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.ledger_entries WHERE id=$1", entry.Id).Scan(&count)
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

		newlyCreatedTransaction, err := repository.CreateTransaction(t.Context(), tx, dto.Transaction{
			Identifier: uuid.New(),
			TransferId: 1,
			CreatedAt:  time.Now(),
			Status:     "pending",
		})
		require.NoError(t, err)

		err = tx.Commit(t.Context())
		require.NoError(t, err)

		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.transactions WHERE id=$1", newlyCreatedTransaction.Id).Scan(&count)
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

		newlyCreatedTransaction, err := repository.CreateTransaction(t.Context(), tx, dto.Transaction{
			Identifier: uuid.New(),
			TransferId: 1,
			CreatedAt:  time.Now(),
			Status:     "pending",
		})
		require.NoError(t, err)

		err = tx.Rollback(t.Context())
		require.NoError(t, err)

		var count int
		err = pool.QueryRow(t.Context(), "SELECT count(*) FROM ledger.transactions WHERE id=$1", newlyCreatedTransaction.Id).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}

func TestGettingAccountBalance(t *testing.T) {
	repository := NewLedgerPostgresRepository(pool)
	t.Run("when transaction is not aborted", func(t *testing.T) {
		tx, err := pool.Begin(t.Context())
		if err != nil {
			panic(err)
		}
		defer tx.Rollback(t.Context())

		accountBalance, err := repository.GetAccountBalance(t.Context(), tx, 1)
		require.NoError(t, err)

		err = tx.Commit(t.Context())
		require.NoError(t, err)

		timee, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-15 12:00:00 +0000")
		require.NoError(t, err)

		assert.Equal(t, int64(1), accountBalance.AccountId)
		assert.Equal(t, float64(100), accountBalance.Availablebalance)
		assert.Equal(t, timee, accountBalance.UpdatedAt)
	})
}

func TestUpdatingAccountBalance(t *testing.T) {
	repository := NewLedgerPostgresRepository(pool)

	t.Run("when transaction is not aborted", func(t *testing.T) {
		tx, err := pool.Begin(t.Context())
		if err != nil {
			panic(err)
		}
		defer tx.Rollback(t.Context())

		timee, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-02-15 12:00:00 +0000")
		require.NoError(t, err)
		repository.UpdateAccountBalance(t.Context(), tx, dto.AccountBalance{
			AccountId:        1,
			Availablebalance: 201,
			UpdatedAt:        timee,
		})

		require.NoError(t, tx.Commit(t.Context()))

		var availableBalance float64
		var updatedAt time.Time
		err = pool.QueryRow(t.Context(), "SELECT available_balance, updated_at FROM ledger.account_balances WHERE account_id = $1", 1).Scan(&availableBalance, &updatedAt)
		require.NoError(t, err)

		assert.Equal(t, float64(201), availableBalance)
		assert.Equal(t, timee, updatedAt)
	})

	t.Run("when transaction is aborted", func(t *testing.T) {
		tx, err := pool.Begin(t.Context())
		if err != nil {
			panic(err)
		}

		timee, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-02-15 12:00:00 +0000")
		require.NoError(t, err)
		repository.UpdateAccountBalance(t.Context(), tx, dto.AccountBalance{
			AccountId:        2,
			Availablebalance: 201,
			UpdatedAt:        timee,
		})

		require.NoError(t, tx.Rollback(t.Context()))

		var availableBalance float64
		var updatedAt time.Time
		err = pool.QueryRow(t.Context(), "SELECT available_balance, updated_at FROM ledger.account_balances WHERE account_id = $1", 2).Scan(&availableBalance, &updatedAt)
		require.NoError(t, err)

		assert.Equal(t, float64(100), availableBalance)

		timeeActual, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-15 12:00:00 +0000")
		require.NoError(t, err)
		assert.Equal(t, timeeActual, updatedAt)
	})
}
