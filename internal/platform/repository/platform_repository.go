package repository

import (
	"context"

	"github.com/deshortone/ledger-system/internal/platform/database_base"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransferPostgresRepository struct {
	database_base.BaseRepo
}

func NewPlatformPostgresRepository(pool *pgxpool.Pool) TransferPostgresRepository {
	if pool == nil {
		panic("pool cannot be nil")
	}
	return TransferPostgresRepository{
		BaseRepo: database_base.NewBaseRepo(pool),
	}
}

func (r TransferPostgresRepository) IsUp(ctx context.Context) error {
	executor := r.GetExecutor(ctx)
	_, err := executor.Exec(ctx, "SELECT 1")
	return err
}
