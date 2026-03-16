package postgres

import (
	"context"
	"strconv"

	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/deshortone/ledger-system/internal/ledger/repository/postgres/ledgerdb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerPostgresRepository struct {
	queries *ledgerdb.Queries
}

func NewLedgerPostgresRepository(pool *pgxpool.Pool) *LedgerPostgresRepository {
	if pool == nil {
		panic("pool cannot be nil")
	}
	return &LedgerPostgresRepository{
		queries: ledgerdb.New(pool),
	}
}

func (r *LedgerPostgresRepository) CreateLedgerEntry(ctx context.Context, tx pgx.Tx, record dto.LedgerEntry) error {
	queries := r.queries.WithTx(tx)
	balance, err := Float64ToNumeric(record.Amount)
	if err != nil {
		return err
	}
	return queries.CreateLedgerEntry(ctx, ledgerdb.CreateLedgerEntryParams{
		Identifier:   record.Identifier,    // yikes
		Identifier_2: record.TransactionId, // yikes_2
		Identifier_3: record.AccountId,     // yikes_3
		Amount:       balance,
		Direction:    ledgerdb.LedgerEntryDirection(record.Direction),
		CreatedAt:    record.CreatedAt,
	})
}

func (r *LedgerPostgresRepository) CreateTransaction(ctx context.Context, tx pgx.Tx, record dto.Transaction) error {
	queries := r.queries.WithTx(tx)
	return queries.CreateTransaction(ctx, ledgerdb.CreateTransactionParams{
		Identifier:   record.Identifier,
		Identifier_2: record.TransferId,
		CreatedAt:    record.CreatedAt,
		Status:       record.Status,
	})
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
