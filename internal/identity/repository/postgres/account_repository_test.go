package postgres

import (
	"testing"
	"time"

	"github.com/deshortone/ledger-system/internal/identity/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanHandleAccounts(t *testing.T) {
	users := initaliseUsers(t)

	r := NewAccountPostgresRepository(pool)
	t.Run("When user owns no accounts", func(t *testing.T) {
		accounts, err := r.GetAccountsOwnedByUser(t.Context(), users[0])
		require.NoError(t, err)
		assert.Empty(t, accounts)
	})

	t.Run("Successfully add and get accounts assigned to user 1 but not user 2", func(t *testing.T) {
		accountsAdded := make([]dto.Account, 0, 2)
		for range 2 {
			account := dto.Account{
				Identifier:     uuid.New(),
				UserIdentifier: users[0].Identifier,
				CreatedAt:      dto.NewCustomTime(time.Now()), // flakey due to timezones?
				AccountType:    "Checking",
				Currency:       "GBP",
				Status:         "active",
			}

			err := r.CreateAccount(t.Context(), account)
			require.NoError(t, err)
			accountsAdded = append(accountsAdded, account)
		}

		accountsReturnedForUser, err := r.GetAccountsOwnedByUser(t.Context(), users[0])
		require.NoError(t, err)
		assert.ElementsMatch(t, accountsAdded, accountsReturnedForUser)

		accounts, err := r.GetAccountsOwnedByUser(t.Context(), users[1])
		require.NoError(t, err)
		assert.Empty(t, accounts)
	})
}

func initaliseUsers(t *testing.T) []dto.User {
	t.Helper()

	users := make([]dto.User, 0, 2)
	for range 2 {
		userToAdd := dto.User{
			Identifier: uuid.New(),
			FirstName:  "the first name",
			LastName:   "the last name",
		}

		r := NewUserPostgresRepository(pool)
		err := r.CreateUser(t.Context(), userToAdd)
		require.NoError(t, err)
		users = append(users, userToAdd)
	}
	return users
}
