package service

import (
	"context"

	"github.com/deshortone/ledger-system/internal/identity/common"
	"github.com/deshortone/ledger-system/internal/identity/domain"
	"github.com/deshortone/ledger-system/internal/identity/dto"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) UserService {
	return UserService{
		userRepository: userRepository,
	}
}

func (s UserService) CreateNewUser(ctx context.Context, firstName, lastName string) (dto.User, error) {
	if len(firstName) == 0 || len(lastName) == 0 {
		return dto.User{}, common.ErrUserDetailsNotFilledIn
	}

	userToCreate := dto.User{
		Identifier: uuid.New(),
		FirstName:  firstName,
		LastName:   lastName,
	}

	err := s.userRepository.CreateUser(ctx, userToCreate)
	if err != nil {
		return dto.User{}, err
	}

	return userToCreate, nil
}
