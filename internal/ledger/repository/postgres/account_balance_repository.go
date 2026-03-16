package postgres

import (
	"context"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/deshortone/ledger-system/internal/ledger/repository/postgres/accountbalance"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountBalancePostgresRepository struct {
	queries *accountbalance.Queries
}

func NewAccountBalancePostgresRepository(pool *pgxpool.Pool) *AccountBalancePostgresRepository {
	if pool == nil {
		panic("pool cannot be nil")
	}
	return &AccountBalancePostgresRepository{
		queries: accountbalance.New(pool),
	}
}

func (r *AccountBalancePostgresRepository) CreateNewAccountBalance(ctx context.Context, accountId uuid.UUID, createdAt time.Time) error {
	balance, err := Float64ToNumeric(0)
	if err != nil {
		return err
	}
	return r.queries.CreateAccountBalance(ctx, accountbalance.CreateAccountBalanceParams{
		Identifier:       accountId,
		AvailableBalance: balance,
		UpdatedAt:        createdAt,
	})
}

func (r *AccountBalancePostgresRepository) GetAccountBalance(ctx context.Context, tx pgx.Tx, accountId uuid.UUID) (dto.AccountBalance, error) {
	queries := r.queries.WithTx(tx)
	accountBalanceRecord, err := queries.GetAccountBalanceAndLock(ctx, accountId)
	if err != nil {
		return dto.AccountBalance{}, err
	}

	accountBalance, err := NumericToFloat64(accountBalanceRecord.AvailableBalance)
	if err != nil {
		return dto.AccountBalance{}, err
	}

	return dto.AccountBalance{
		AccountId:        accountId,
		Availablebalance: accountBalance,
		UpdatedAt:        accountBalanceRecord.UpdatedAt,
	}, nil
}

func (r *AccountBalancePostgresRepository) UpdateAccountBalance(ctx context.Context, tx pgx.Tx, record dto.AccountBalance) error {
	queries := r.queries.WithTx(tx)
	balance, err := Float64ToNumeric(record.Availablebalance)
	if err != nil {
		return err
	}
	return queries.UpdateAccountBalance(ctx, accountbalance.UpdateAccountBalanceParams{
		Identifier:       record.AccountId,
		AvailableBalance: balance,
		UpdatedAt:        record.UpdatedAt,
	})
}
