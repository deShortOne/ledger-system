package service

import (
	"testing"

	"github.com/deshortone/ledger-system/internal/identity/common"
	"github.com/deshortone/ledger-system/internal/identity/repository/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	service := NewUserService(memory.NewUserInMemoryRepository())

	t.Run("fail as firstname is empty", func(t *testing.T) {
		_, err := service.CreateNewUser(t.Context(), "", "last name")
		assert.ErrorIs(t, common.ErrUserDetailsNotFilledIn, err)
	})

	t.Run("fail as lastname is empty", func(t *testing.T) {
		_, err := service.CreateNewUser(t.Context(), "first name", "")
		assert.ErrorIs(t, common.ErrUserDetailsNotFilledIn, err)
	})

	t.Run("Succeed", func(t *testing.T) {
		userCreated, err := service.CreateNewUser(t.Context(), "first name", "last name")
		require.NoError(t, err)
		assert.NotNil(t, userCreated.Identifier)
		assert.Equal(t, "first name", userCreated.FirstName)
		assert.Equal(t, "last name", userCreated.LastName)
	})
}
