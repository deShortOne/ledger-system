package service

import (
	"testing"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/contracts"
	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/deshortone/ledger-system/internal/ledger/repository/memory"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddToLedger(t *testing.T) {
	t.Run("when the transaction is going to be successful", func(t *testing.T) {
		t.Parallel()
		var tx pgx.Tx

		timeAccount1, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-02-15 12:00:00 +0000")
		require.NoError(t, err)
		repository := memory.NewLedgerInMemoryRepository()
		err = repository.CreateAccountBalance(t.Context(), tx, dto.AccountBalance{
			AccountId:        1,
			Availablebalance: 100,
			UpdatedAt:        timeAccount1,
		})
		require.NoError(t, err)
		err = repository.CreateAccountBalance(t.Context(), tx, dto.AccountBalance{
			AccountId:        2,
			Availablebalance: 100,
			UpdatedAt:        timeAccount1,
		})
		require.NoError(t, err)

		service := NewLedgerService(repository)

		timeOfTransfer, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-15 12:00:00 +0000")
		err = service.AddToLedger(t.Context(), tx, contracts.AddToLedgerRequest{
			TransferId: 1,
			CreatedAt:  timeOfTransfer,
			Entries: []contracts.LedgerEntries{
				{
					AccountId: 1,
					Amount:    10,
					Direction: contracts.CREDIT,
				},
				{
					AccountId: 2,
					Amount:    10,
					Direction: contracts.DEBIT,
				},
			},
		})

		assert.NoError(t, err)

		account1, err := repository.GetAccountBalance(t.Context(), tx, 1)
		require.NoError(t, err)
		account2, err := repository.GetAccountBalance(t.Context(), tx, 2)
		require.NoError(t, err)
		assert.Equal(t, float64(110), account1.Availablebalance)
		assert.Equal(t, timeOfTransfer, account1.UpdatedAt)
		assert.Equal(t, float64(90), account2.Availablebalance)
		assert.Equal(t, timeOfTransfer, account2.UpdatedAt)
	})

	t.Run("when the double entry is violated", func(t *testing.T) {
		t.Parallel()
		var tx pgx.Tx

		timeAccount1, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-02-15 12:00:00 +0000")
		require.NoError(t, err)
		repository := memory.NewLedgerInMemoryRepository()
		err = repository.CreateAccountBalance(t.Context(), tx, dto.AccountBalance{
			AccountId:        1,
			Availablebalance: 100,
			UpdatedAt:        timeAccount1,
		})
		require.NoError(t, err)
		err = repository.CreateAccountBalance(t.Context(), tx, dto.AccountBalance{
			AccountId:        2,
			Availablebalance: 100,
			UpdatedAt:        timeAccount1,
		})
		require.NoError(t, err)

		service := NewLedgerService(repository)

		timeOfTransfer, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-15 12:00:00 +0000")
		err = service.AddToLedger(t.Context(), tx, contracts.AddToLedgerRequest{
			TransferId: 1,
			CreatedAt:  timeOfTransfer,
			Entries: []contracts.LedgerEntries{
				{
					AccountId: 1,
					Amount:    10,
					Direction: contracts.CREDIT,
				},
				{
					AccountId: 2,
					Amount:    20,
					Direction: contracts.DEBIT,
				},
			},
		})

		assert.ErrorIs(t, err, contracts.ErrDoubleEntryViolated)
	})

	t.Run("when the the debiting account doesn't have enough money", func(t *testing.T) {
		t.Parallel()
		var tx pgx.Tx

		timeAccount1, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-02-15 12:00:00 +0000")
		require.NoError(t, err)
		repository := memory.NewLedgerInMemoryRepository()
		err = repository.CreateAccountBalance(t.Context(), tx, dto.AccountBalance{
			AccountId:        1,
			Availablebalance: 100,
			UpdatedAt:        timeAccount1,
		})
		require.NoError(t, err)
		err = repository.CreateAccountBalance(t.Context(), tx, dto.AccountBalance{
			AccountId:        2,
			Availablebalance: 5,
			UpdatedAt:        timeAccount1,
		})
		require.NoError(t, err)

		service := NewLedgerService(repository)

		timeOfTransfer, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-15 12:00:00 +0000")
		err = service.AddToLedger(t.Context(), tx, contracts.AddToLedgerRequest{
			TransferId: 1,
			CreatedAt:  timeOfTransfer,
			Entries: []contracts.LedgerEntries{
				{
					AccountId: 1,
					Amount:    10,
					Direction: contracts.CREDIT,
				},
				{
					AccountId: 2,
					Amount:    10,
					Direction: contracts.DEBIT,
				},
			},
		})

		assert.ErrorIs(t, err, contracts.ErrOneOfTheAccountsDoNotHaveEnoughMoney)
	})
}
