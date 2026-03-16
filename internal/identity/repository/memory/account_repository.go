package memory

import (
	"context"

	"github.com/deshortone/ledger-system/internal/identity/dto"
)

type AccountInMemoryRepository struct {
	accounts []dto.Account
}

func NewAccountInMemoryRepository() *AccountInMemoryRepository {
	return &AccountInMemoryRepository{
		accounts: []dto.Account{},
	}
}

func (r *AccountInMemoryRepository) CreateAccount(ctx context.Context, account dto.Account) (dto.Account, error) {
	account.Id = int64(len(r.accounts) + 1)
	r.accounts = append(r.accounts, account)
	return account, nil
}

func (r *AccountInMemoryRepository) GetAccountsOwnedByUser(ctx context.Context, user dto.User) ([]dto.Account, error) {
	accounts := []dto.Account{}
	for _, account := range r.accounts {
		if account.UserId == user.Id {
			accounts = append(accounts, account)
		}
	}

	return accounts, nil
}
