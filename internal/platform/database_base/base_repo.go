package database_base

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepo struct {
	pool *pgxpool.Pool
}

func NewBaseRepo(pool *pgxpool.Pool) BaseRepo {
	return BaseRepo{
		pool: pool,
	}
}

func (r *BaseRepo) GetExecutor(ctx context.Context) DBTX {
	if tx, ok := getTx(ctx); ok {
		return tx
	}
	return r.pool
}
