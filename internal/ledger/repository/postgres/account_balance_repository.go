package postgres

import (
	"context"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/deshortone/ledger-system/internal/ledger/repository/postgres/accountbalance"
	"github.com/deshortone/ledger-system/internal/platform/database_base"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountBalancePostgresRepository struct {
	database_base.BaseRepo
}

func NewAccountBalancePostgresRepository(pool *pgxpool.Pool) *AccountBalancePostgresRepository {
	if pool == nil {
		panic("pool cannot be nil")
	}
	return &AccountBalancePostgresRepository{
		BaseRepo: database_base.NewBaseRepo(pool),
	}
}

func (r *AccountBalancePostgresRepository) CreateNewAccountBalance(ctx context.Context, accountId uuid.UUID, createdAt time.Time) error {
	executor := r.GetExecutor(ctx)
	queries := accountbalance.New(executor)

	balance, err := Float64ToNumeric(0)
	if err != nil {
		return err
	}
	return queries.CreateAccountBalance(ctx, accountbalance.CreateAccountBalanceParams{
		Identifier:       accountId,
		AvailableBalance: balance,
		UpdatedAt:        createdAt,
	})
}

func (r *AccountBalancePostgresRepository) GetAccountBalance(ctx context.Context, accountId uuid.UUID) (dto.AccountBalance, error) {
	executor := r.GetExecutor(ctx)
	queries := accountbalance.New(executor)

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

func (r *AccountBalancePostgresRepository) UpdateAccountBalance(ctx context.Context, record dto.AccountBalance) error {
	executor := r.GetExecutor(ctx)
	queries := accountbalance.New(executor)

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
