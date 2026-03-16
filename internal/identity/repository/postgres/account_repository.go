package postgres

import (
	"context"

	"github.com/deshortone/ledger-system/internal/identity/dto"
	accountdb "github.com/deshortone/ledger-system/internal/identity/repository/postgres/account_db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountPostgresRepository struct {
	queries *accountdb.Queries
}

func NewAccountPostgresRepository(pool *pgxpool.Pool) *AccountPostgresRepository {
	if pool == nil {
		panic("pool cannot be nil")
	}

	return &AccountPostgresRepository{
		queries: accountdb.New(pool),
	}
}

func (r *AccountPostgresRepository) CreateAccount(ctx context.Context, account dto.Account) error {
	err := r.queries.CreateAccount(ctx, accountdb.CreateAccountParams{
		Identifier:   account.Identifier,
		Identifier_2: account.UserIdentifier,
		CreatedAt:    account.CreatedAt.Time,
		AccountType:  account.AccountType,
		Currency:     account.Currency,
		Status:       account.Status,
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *AccountPostgresRepository) GetAccountsOwnedByUser(ctx context.Context, user dto.User) ([]dto.Account, error) {
	accounts, err := r.queries.GetAccountsOwnedByUser(ctx, user.Identifier)
	if err != nil {
		return []dto.Account{}, err
	}

	accountsResponse := []dto.Account{}
	for _, account := range accounts {
		accountsResponse = append(accountsResponse, dto.Account{
			Identifier:     account.AccountIdentifier,
			UserIdentifier: user.Identifier,
			CreatedAt:      dto.NewCustomTime(account.CreatedAt),
			AccountType:    account.AccountType,
			Currency:       account.Currency,
			Status:         account.Status,
		})
	}

	return accountsResponse, nil
}
