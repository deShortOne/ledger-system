package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/deshortone/ledger-system/internal/platform/database_base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGettingAccountBalance(t *testing.T) {
	repository := NewAccountBalancePostgresRepository(pool)
	t.Run("when transaction is not aborted", func(t *testing.T) {
		accountBalance, err := repository.GetAccountBalance(t.Context(), account1Id)
		require.NoError(t, err)

		timee, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-15 12:00:00 +0000")
		require.NoError(t, err)

		assert.Equal(t, account1Id, accountBalance.AccountId)
		assert.Equal(t, float64(100), accountBalance.Availablebalance)
		assert.Equal(t, timee, accountBalance.UpdatedAt)
	})
}

func TestUpdatingAccountBalance(t *testing.T) {
	repository := NewAccountBalancePostgresRepository(pool)

	t.Run("when transaction is not aborted", func(t *testing.T) {
		timee, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-02-15 12:00:00 +0000")
		require.NoError(t, err)
		repository.UpdateAccountBalance(t.Context(), dto.AccountBalance{
			AccountId:        account1Id,
			Availablebalance: 201,
			UpdatedAt:        timee,
		})

		var availableBalance float64
		var updatedAt time.Time
		err = pool.QueryRow(t.Context(), "SELECT available_balance, updated_at FROM ledger.account_balances WHERE account_id = $1", 1).Scan(&availableBalance, &updatedAt)
		require.NoError(t, err)

		assert.Equal(t, float64(201), availableBalance)
		assert.Equal(t, timee, updatedAt)
	})

	t.Run("when transaction is aborted", func(t *testing.T) {
		timee, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-02-15 12:00:00 +0000")
		require.NoError(t, err)
		uow := database_base.NewPgUnitOfWork(pool)
		err = uow.Do(t.Context(), func(ctx context.Context) error {
			err = repository.UpdateAccountBalance(ctx, dto.AccountBalance{
				AccountId:        account2Id,
				Availablebalance: 201,
				UpdatedAt:        timee,
			})
			require.NoError(t, err)

			return errors.New("throw something")
		})
		require.Error(t, err)

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
