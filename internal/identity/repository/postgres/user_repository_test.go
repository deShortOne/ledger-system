package postgres

import (
	"testing"

	"github.com/deshortone/ledger-system/internal/identity/common"
	"github.com/deshortone/ledger-system/internal/identity/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanHandleUsers(t *testing.T) {
	userToAdd := dto.User{
		Identifier: uuid.New(),
		FirstName:  "the first name",
		LastName:   "the last name",
	}

	r := NewUserPostgresRepository(pool)
	t.Run("Throw error when getting non existing user", func(t *testing.T) {
		_, err := r.GetUser(t.Context(), userToAdd.Identifier)

		assert.ErrorIs(t, common.ErrUserIdentifierNotFound, err)
	})

	t.Run("Successfully add and get new user", func(t *testing.T) {
		err := r.CreateUser(t.Context(), userToAdd)
		require.NoError(t, err)

		user, err := r.GetUser(t.Context(), userToAdd.Identifier)
		assert.NoError(t, err)
		assert.Equal(t, userToAdd.Identifier, user.Identifier)
		assert.Equal(t, userToAdd.FirstName, user.FirstName)
		assert.Equal(t, userToAdd.LastName, user.LastName)
		assert.Equal(t, userToAdd, user)
	})
}
