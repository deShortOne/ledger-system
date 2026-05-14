package postgres

import (
	"context"
	"strconv"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/deshortone/ledger-system/internal/ledger/repository/postgres/ledgerdb"
	"github.com/deshortone/ledger-system/internal/platform/database_base"
	"github.com/deshortone/ledger-system/pkg/failure"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerPostgresRepository struct {
	database_base.BaseRepo
}

func NewLedgerPostgresRepository(pool *pgxpool.Pool) *LedgerPostgresRepository {
	if pool == nil {
		panic("pool cannot be nil")
	}
	return &LedgerPostgresRepository{
		BaseRepo: database_base.NewBaseRepo(pool),
	}
}

func (r *LedgerPostgresRepository) CreateLedgerEntry(ctx context.Context, record dto.LedgerEntry) error {
	executor := r.GetExecutor(ctx)
	queries := ledgerdb.New(executor)

	balance, err := Float64ToNumeric(record.Amount)
	if err != nil {
		return failure.NewFailure(
			failure.ConversionError,
			failure.GeneralFailure,
			err,
			"Failed to convert ledger entry amount",
		)
	}

	err = queries.CreateLedgerEntry(ctx, ledgerdb.CreateLedgerEntryParams{
		Identifier:   record.Identifier,    // yikes
		Identifier_2: record.TransactionId, // yikes_2
		Identifier_3: record.AccountId,     // yikes_3
		Amount:       balance,
		Direction:    ledgerdb.LedgerEntryDirection(record.Direction),
		CreatedAt:    record.CreatedAt,
	})
	if err != nil {
		return failure.NewFailure(
			failure.UnknownRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to create ledger entry",
		)
	}

	return nil
}

func (r *LedgerPostgresRepository) CreateTransaction(ctx context.Context, record dto.Transaction) error {
	executor := r.GetExecutor(ctx)
	queries := ledgerdb.New(executor)

	err := queries.CreateTransaction(ctx, ledgerdb.CreateTransactionParams{
		Identifier:   record.Identifier,
		Identifier_2: record.TransferId,
		CreatedAt:    record.CreatedAt,
		Status:       record.Status,
	})
	if err != nil {
		return failure.NewFailure(
			failure.UnknownRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to create transaction",
		)
	}

	return nil
}

func NumericToFloat64(n pgtype.Numeric) (float64, error) {
	f, err := n.Float64Value()
	if err != nil {
		return 0, err
	}
	return f.Float64, nil
}

func Float64ToNumeric(f float64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	err := n.Scan(strconv.FormatFloat(f, 'f', 6, 64))
	return n, err
}
