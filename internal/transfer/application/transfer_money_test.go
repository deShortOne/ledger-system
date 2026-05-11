package application

import (
	"testing"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/contracts"
	"github.com/deshortone/ledger-system/internal/ledger/dto"
	account_memory "github.com/deshortone/ledger-system/internal/ledger/repository/memory"
	ledger_memory "github.com/deshortone/ledger-system/internal/ledger/repository/memory"
	ledger_service "github.com/deshortone/ledger-system/internal/ledger/service"
	"github.com/deshortone/ledger-system/internal/platform/database_base"
	transfer_memory "github.com/deshortone/ledger-system/internal/transfer/repository/memory"
	transfer_service "github.com/deshortone/ledger-system/internal/transfer/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransferringMoney(t *testing.T) {
	t.Run("when everything goes smoothly", func(t *testing.T) {
		transferRepo := transfer_memory.NewTransferInMemoryRepository()
		transferService := transfer_service.NewTransferService(transferRepo)

		accountRepo := account_memory.NewAccountBalanceInMemoryRepository()
		ledgerRepo := ledger_memory.NewLedgerInMemoryRepository()
		ledgerService := ledger_service.NewLedgerService(ledgerRepo, accountRepo)

		uow := database_base.FakeUOW{}

		transferMoneyApplication := NewTransferMoneyBetweenAccounts(transferService, ledgerService, uow)

		fromAccountId := uuid.New()
		toAccountId := uuid.New()
		amount := float64(523)
		accountRepo.CreateNewAccountBalance(t.Context(), fromAccountId, time.Now())
		accountRepo.UpdateAccountBalance(t.Context(), dto.AccountBalance{
			AccountId:        fromAccountId,
			Availablebalance: 600,
			UpdatedAt:        time.Now(),
		})
		accountRepo.CreateNewAccountBalance(t.Context(), toAccountId, time.Now())

		transferId, err := transferMoneyApplication.TransferMoney(t.Context(), fromAccountId, toAccountId, amount)
		require.NotEqual(t, uuid.Nil, transferId)
		require.NoError(t, err)

		acc1, err := accountRepo.GetAccountBalance(t.Context(), fromAccountId)
		require.NoError(t, err)
		assert.Equal(t, float64(77), acc1.Availablebalance)

		acc2, err := accountRepo.GetAccountBalance(t.Context(), toAccountId)
		require.NoError(t, err)
		assert.Equal(t, float64(523), acc2.Availablebalance)

		require.Equal(t, 1, len(transferRepo.TransferRequests))
		transferRequestAcc := transferRepo.TransferRequests[0]
		assert.Equal(t, fromAccountId, transferRequestAcc.FromAccountId)
		assert.Equal(t, toAccountId, transferRequestAcc.ToAccountId)
		assert.Equal(t, float64(523), transferRequestAcc.Amount)
		assert.Equal(t, "pending", transferRequestAcc.Status)

		require.Equal(t, 1, len(transferRepo.StatusUpdates))
		statusUpdateAcc := transferRepo.StatusUpdates[0]
		assert.Equal(t, transferRequestAcc.Identifier, statusUpdateAcc.Id)
		assert.Equal(t, "Success", statusUpdateAcc.Status)
		assert.Equal(t, false, statusUpdateAcc.WasFailureUpdated)

		require.Equal(t, 1, len(transferRepo.Transfers))
		transferAcc := transferRepo.Transfers[0]
		assert.Equal(t, transferRequestAcc.Identifier, transferAcc.TransferRequestId)

		require.Equal(t, 1, len(ledgerRepo.Transactions))
		transaction := ledgerRepo.Transactions[0]
		assert.Equal(t, "posted", transaction.Status)
		assert.Equal(t, transferAcc.Identifier, transaction.TransferId)

		require.Equal(t, 2, len(ledgerRepo.LedgerEntries))
		ledgerEntry1 := ledgerRepo.LedgerEntries[0]
		assert.Equal(t, fromAccountId, ledgerEntry1.AccountId)
		assert.Equal(t, float64(523), ledgerEntry1.Amount)
		assert.Equal(t, contracts.LedgerDirection("DEBIT"), ledgerEntry1.Direction)
		assert.Equal(t, transaction.Identifier, ledgerEntry1.TransactionId)

		ledgerEntry2 := ledgerRepo.LedgerEntries[1]
		assert.Equal(t, toAccountId, ledgerEntry2.AccountId)
		assert.Equal(t, float64(523), ledgerEntry2.Amount)
		assert.Equal(t, contracts.LedgerDirection("CREDIT"), ledgerEntry2.Direction)
		assert.Equal(t, transaction.Identifier, ledgerEntry2.TransactionId)
	})

	t.Run("when there is an error during the unit of work", func(t *testing.T) {
		transferRepo := transfer_memory.NewTransferInMemoryRepository()
		transferService := transfer_service.NewTransferService(transferRepo)

		accountRepo := account_memory.NewAccountBalanceInMemoryRepository()
		ledgerRepo := ledger_memory.NewLedgerInMemoryRepository()
		ledgerService := ledger_service.NewLedgerService(ledgerRepo, accountRepo)

		uow := database_base.FakeUOW{}

		transferMoneyApplication := NewTransferMoneyBetweenAccounts(transferService, ledgerService, uow)

		fromAccountId := uuid.New()
		toAccountId := uuid.New()
		amount := float64(523)
		accountRepo.CreateNewAccountBalance(t.Context(), fromAccountId, time.Now())
		accountRepo.UpdateAccountBalance(t.Context(), dto.AccountBalance{
			AccountId:        fromAccountId,
			Availablebalance: 100,
			UpdatedAt:        time.Now(),
		})
		accountRepo.CreateNewAccountBalance(t.Context(), toAccountId, time.Now())

		transferId, err := transferMoneyApplication.TransferMoney(t.Context(), fromAccountId, toAccountId, amount)
		require.NotEqual(t, uuid.Nil, transferId) // as long as the transfer exists in the db, we should always return the id for it
		require.NoError(t, err)

		acc1, err := accountRepo.GetAccountBalance(t.Context(), fromAccountId)
		require.NoError(t, err)
		assert.Equal(t, float64(100), acc1.Availablebalance)

		acc2, err := accountRepo.GetAccountBalance(t.Context(), toAccountId)
		require.NoError(t, err)
		assert.Equal(t, float64(0), acc2.Availablebalance)

		require.Equal(t, 1, len(transferRepo.TransferRequests))
		transferRequestAcc := transferRepo.TransferRequests[0]
		assert.Equal(t, fromAccountId, transferRequestAcc.FromAccountId)
		assert.Equal(t, toAccountId, transferRequestAcc.ToAccountId)
		assert.Equal(t, float64(523), transferRequestAcc.Amount)
		assert.Equal(t, "pending", transferRequestAcc.Status)

		require.Equal(t, 1, len(transferRepo.StatusUpdates))
		statusUpdateAcc := transferRepo.StatusUpdates[0]
		assert.Equal(t, transferRequestAcc.Identifier, statusUpdateAcc.Id)
		assert.Equal(t, "Failed", statusUpdateAcc.Status)
		assert.Equal(t, true, statusUpdateAcc.WasFailureUpdated)

		require.Equal(t, 1, len(transferRepo.Transfers))
		transferAcc := transferRepo.Transfers[0]
		assert.Equal(t, transferRequestAcc.Identifier, transferAcc.TransferRequestId)

		require.Equal(t, 1, len(ledgerRepo.Transactions))
		transaction := ledgerRepo.Transactions[0]
		assert.Equal(t, "posted", transaction.Status)
		assert.Equal(t, transferAcc.Identifier, transaction.TransferId)

		require.Equal(t, 0, len(ledgerRepo.LedgerEntries))
	})
}
