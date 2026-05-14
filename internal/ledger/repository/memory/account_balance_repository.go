package memory

import (
	"context"
	"errors"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/deshortone/ledger-system/pkg/failure"
	"github.com/google/uuid"
)

type AccountBalanceInMemoryRepository struct {
	accountBalances []dto.AccountBalance
}

func NewAccountBalanceInMemoryRepository() *AccountBalanceInMemoryRepository {
	return &AccountBalanceInMemoryRepository{
		accountBalances: []dto.AccountBalance{},
	}
}

func (r *AccountBalanceInMemoryRepository) GetAccountBalance(ctx context.Context, accountId uuid.UUID) (dto.AccountBalance, error) {
	for _, accountBalance := range r.accountBalances {
		if accountBalance.AccountId == accountId {
			return accountBalance, nil
		}
	}

	return dto.AccountBalance{}, failure.NewFailure(
		failure.LedgerNotFound,
		failure.NotFound,
		errors.New("account balance not found"),
		"No account balance could be found for the requested identifier",
	)
}

func (r *AccountBalanceInMemoryRepository) UpdateAccountBalance(ctx context.Context, record dto.AccountBalance) error {
	for i, accountBalance := range r.accountBalances {
		if accountBalance.AccountId == record.AccountId {
			r.accountBalances[i].Availablebalance = record.Availablebalance
			r.accountBalances[i].UpdatedAt = record.UpdatedAt
			return nil
		}
	}

	return failure.NewFailure(
		failure.LedgerNotFound,
		failure.NotFound,
		errors.New("account balance not found"),
		"The account balance for the requested identifier was not available to update",
	)
}

func (r *AccountBalanceInMemoryRepository) CreateNewAccountBalance(ctx context.Context, accountId uuid.UUID, createdAt time.Time) error {
	r.accountBalances = append(r.accountBalances, dto.AccountBalance{
		AccountId:        accountId,
		Availablebalance: 0,
		UpdatedAt:        createdAt,
	})
	return nil
}
