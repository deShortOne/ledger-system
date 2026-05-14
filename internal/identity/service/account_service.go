package service

import (
	"context"
	"errors"

	"github.com/deshortone/ledger-system/internal/identity/common"
	"github.com/deshortone/ledger-system/internal/identity/domain"
	"github.com/deshortone/ledger-system/internal/identity/dto"
	"github.com/deshortone/ledger-system/pkg/failure"
	"github.com/google/uuid"
)

type AccountService struct {
	accountRepository domain.AccountRepository
	userRepository    domain.UserRepository
}

func NewAccountService(
	accountRepository domain.AccountRepository,
	userRepository domain.UserRepository,
) AccountService {
	return AccountService{
		accountRepository: accountRepository,
		userRepository:    userRepository,
	}
}

func (s AccountService) AddAccountToUser(ctx context.Context,
	userIdentifier uuid.UUID,
	accountType, currency string,
) (dto.Account, error) {
	user, err := s.userRepository.GetUser(ctx, userIdentifier)
	if err != nil {
		if errors.Is(err, common.ErrUserIdentifierNotFound) {
			return dto.Account{}, failure.NewFailure(
				failure.UserNotFound,
				failure.NotFound,
				err,
				"Specified user identifier does not exist",
			)
		}
		return dto.Account{}, failure.NewFailure(
			failure.IdentityRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to retrieve user before creating account",
		)
	}

	account := dto.Account{
		Identifier:     uuid.New(),
		UserIdentifier: user.Identifier,
		CreatedAt:      dto.NewCustomTimeNow(),
		AccountType:    accountType,
		Currency:       currency,
		Status:         "available",
	}
	if err = s.accountRepository.CreateAccount(ctx, account); err != nil {
		return dto.Account{}, failure.NewFailure(
			failure.IdentityRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to create account for user",
		)
	}

	return account, nil
}

func (s AccountService) GetAccountsOwnedByUser(ctx context.Context, identifier uuid.UUID) ([]dto.Account, error) {
	user, err := s.userRepository.GetUser(ctx, identifier)
	if err != nil {
		return []dto.Account{}, err
	}

	return s.accountRepository.GetAccountsOwnedByUser(ctx, user)
}
