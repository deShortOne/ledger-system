package application

import (
	"context"

	"github.com/deshortone/ledger-system/internal/identity/domain"
	"github.com/deshortone/ledger-system/internal/identity/dto"
	"github.com/google/uuid"
)

type CreateNewAccountApplication struct {
	accountService        domain.AccountCreator
	accountBalanceCreator domain.AccountBalanceCreator
}

func NewCreateNewAccountApplication(
	accountService domain.AccountCreator,
	accountBalanceCreator domain.AccountBalanceCreator,
) *CreateNewAccountApplication {
	return &CreateNewAccountApplication{
		accountService:        accountService,
		accountBalanceCreator: accountBalanceCreator,
	}
}

func (a *CreateNewAccountApplication) AddAccountToUser(ctx context.Context, userIdentifier uuid.UUID, accountType, currency string) (dto.Account, error) {
	accountCreated, err := a.accountService.AddAccountToUser(ctx, userIdentifier, accountType, currency)
	if err != nil {
		return dto.Account{}, err
	}

	err = a.accountBalanceCreator.CreateNewAccount(ctx, accountCreated.Identifier, accountCreated.CreatedAt.Time)
	if err != nil {
		return dto.Account{}, err
	}
	return accountCreated, nil
}
