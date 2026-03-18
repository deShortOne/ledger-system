package database_base

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type PgUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewPgUnitOfWork(pool *pgxpool.Pool) *PgUnitOfWork {
	return &PgUnitOfWork{pool: pool}
}

func (u *PgUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}

	ctx = withTx(ctx, tx)

	err = fn(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
