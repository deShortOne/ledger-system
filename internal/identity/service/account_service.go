package service

import (
	"context"

	"github.com/deshortone/ledger-system/internal/identity/domain"
	"github.com/deshortone/ledger-system/internal/identity/dto"
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
	accountToCreate dto.Account,
) (dto.Account, error) {
	user, err := s.userRepository.GetUser(ctx, userIdentifier)
	if err != nil {
		return dto.Account{}, err
	}

	return s.accountRepository.CreateAccount(ctx, dto.Account{
		Identifier:  uuid.New(),
		UserId:      user.Id,
		CreatedAt:   dto.NewCustomTimeNow(),
		AccountType: accountToCreate.AccountType,
		Currency:    accountToCreate.Currency,
		Status:      "available",
	})
}

func (s AccountService) GetAccountsOwnedByUser(ctx context.Context, identifier uuid.UUID) ([]dto.Account, error) {
	user, err := s.userRepository.GetUser(ctx, identifier)
	if err != nil {
		return []dto.Account{}, err
	}

	return s.accountRepository.GetAccountsOwnedByUser(ctx, user)
}
