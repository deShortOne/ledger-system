package postgres

import (
	"context"
	"time"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/deshortone/ledger-system/internal/ledger/repository/postgres/accountbalance"
	"github.com/deshortone/ledger-system/internal/platform/database_base"
	"github.com/deshortone/ledger-system/pkg/failure"
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
		return failure.NewFailure(
			failure.ConversionError,
			failure.GeneralFailure,
			err,
			"Failed to convert initial account balance",
		)
	}

	err = queries.CreateAccountBalance(ctx, accountbalance.CreateAccountBalanceParams{
		Identifier:       accountId,
		AvailableBalance: balance,
		UpdatedAt:        createdAt,
	})
	if err != nil {
		return failure.NewFailure(
			failure.UnknownRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to create account balance",
		)
	}

	return nil
}

func (r *AccountBalancePostgresRepository) GetAccountBalance(ctx context.Context, accountId uuid.UUID) (dto.AccountBalance, error) {
	executor := r.GetExecutor(ctx)
	queries := accountbalance.New(executor)

	accountBalanceRecord, err := queries.GetAccountBalanceAndLock(ctx, accountId)
	if err != nil {
		return dto.AccountBalance{}, failure.NewFailure(
			failure.UnknownRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to get account balance",
		)
	}

	accountBalance, err := NumericToFloat64(accountBalanceRecord.AvailableBalance)
	if err != nil {
		return dto.AccountBalance{}, failure.NewFailure(
			failure.ConversionError,
			failure.GeneralFailure,
			err,
			"Failed to convert account balance",
		)
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
		return failure.NewFailure(
			failure.ConversionError,
			failure.GeneralFailure,
			err,
			"Failed to convert account balance",
		)
	}

	err = queries.UpdateAccountBalance(ctx, accountbalance.UpdateAccountBalanceParams{
		Identifier:       record.AccountId,
		AvailableBalance: balance,
		UpdatedAt:        record.UpdatedAt,
	})
	if err != nil {
		return failure.NewFailure(
			failure.UnknownRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to update account balance",
		)
	}

	return nil
}
