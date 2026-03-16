package memory

import (
	"context"

	"github.com/deshortone/ledger-system/internal/identity/common"
	"github.com/deshortone/ledger-system/internal/identity/dto"
	"github.com/google/uuid"
)

type UserInMemoryRepository struct {
	users []dto.User
}

func NewUserInMemoryRepository() *UserInMemoryRepository {
	return &UserInMemoryRepository{
		users: []dto.User{},
	}
}
func (r *UserInMemoryRepository) CreateUser(ctx context.Context, user dto.User) (dto.User, error) {
	user.Id = int64(len(r.users) + 1)
	r.users = append(r.users, user)
	return user, nil
}

func (r *UserInMemoryRepository) GetUser(ctx context.Context, identifier uuid.UUID) (dto.User, error) {
	for _, user := range r.users {
		if user.Identifier == identifier {
			return user, nil
		}
	}
	return dto.User{}, common.ErrUserIdentifierNotFound
}
