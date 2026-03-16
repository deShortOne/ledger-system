package service

import (
	"testing"

	"github.com/deshortone/ledger-system/internal/identity/common"
	"github.com/deshortone/ledger-system/internal/identity/dto"
	"github.com/deshortone/ledger-system/internal/identity/repository/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	accountService, userIdentifiers := setupAccountServiceTest(t)

	t.Run("when user doesn't exist", func(t *testing.T) {
		badIdentifier, err := uuid.Parse("44e35ed1-0a73-4cc1-b302-925b50a6405e")
		require.NoError(t, err)
		_, err = accountService.AddAccountToUser(t.Context(), badIdentifier, dto.Account{
			AccountType: "",
			Currency:    "",
		})

		assert.ErrorIs(t, common.ErrUserIdentifierNotFound, err)
	})

	t.Run("successfully adding account and retrieving", func(t *testing.T) {
		account, err := accountService.AddAccountToUser(t.Context(), userIdentifiers[0], dto.Account{
			AccountType: "checking",
			Currency:    "GBP",
		})

		require.NoError(t, err)
		assert.Equal(t, "checking", account.AccountType)
		assert.Equal(t, "GBP", account.Currency)

		_, err = accountService.GetAccountsOwnedByUser(t.Context(), userIdentifiers[0])
		require.NoError(t, err)
	})
}

func setupAccountServiceTest(t *testing.T) (AccountService, []uuid.UUID) {
	t.Helper()

	userRepo := memory.NewUserInMemoryRepository()
	userIdentifiers := []uuid.UUID{}
	manualIdentifiers := []string{
		"8be4df61-93ca-11d2-aa0d-00e098032b8c",
		"6052e192-f75a-4c62-a8c9-ec237db718de",
	}
	for i, manualIdentifier := range manualIdentifiers {
		userIdentifier, err := uuid.Parse(manualIdentifier)
		require.NoError(t, err)
		userIdentifiers = append(userIdentifiers, userIdentifier)
		_, err = userRepo.CreateUser(t.Context(), dto.User{
			Id:         int64(i + 1),
			Identifier: userIdentifier,
			FirstName:  "first name",
			LastName:   "last name",
		})
		require.NoError(t, err)
	}

	accountService := NewAccountService(
		memory.NewAccountInMemoryRepository(),
		userRepo,
	)

	return accountService, userIdentifiers
}
