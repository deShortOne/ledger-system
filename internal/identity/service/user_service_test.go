package service

import (
	"context"
	"errors"
	"testing"

	"github.com/deshortone/ledger-system/internal/identity/common"
	"github.com/deshortone/ledger-system/internal/identity/dto"
	"github.com/deshortone/ledger-system/internal/identity/repository/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingUserRepository struct{}

func (failingUserRepository) CreateUser(ctx context.Context, user dto.User) error {
	return errors.New("create user failed")
}

func (failingUserRepository) GetUser(ctx context.Context, identifier uuid.UUID) (dto.User, error) {
	return dto.User{}, nil
}

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

	t.Run("returns repository error", func(t *testing.T) {
		service := NewUserService(failingUserRepository{})

		_, err := service.CreateNewUser(t.Context(), "first name", "last name")
		require.Error(t, err)
		assert.EqualError(t, err, "create user failed")
	})
}
