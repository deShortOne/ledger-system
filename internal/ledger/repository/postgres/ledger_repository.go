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

func (r *LedgerPostgresRepository) CreateLedgerEntry(ctx context.Context, tx pgx.Tx, record dto.LedgerEntry) (dto.LedgerEntry, error) {
	queries := r.queries.WithTx(tx)
	balance, err := Float64ToNumeric(record.Amount)
	if err != nil {
		return dto.LedgerEntry{}, err
	}
	ledgerId, err := queries.CreateLedgerEntry(ctx, ledgerdb.CreateLedgerEntryParams{
		Identifier:    record.Identifier,
		TransactionID: record.TransactionId,
		AccountID:     record.AccountId,
		Amount:        balance,
		Direction:     ledgerdb.LedgerEntryDirection(record.Direction),
		CreatedAt:     record.CreatedAt,
	})
	if err != nil {
		return dto.LedgerEntry{}, err
	}

	record.Id = ledgerId
	return record, nil
}

func (r *LedgerPostgresRepository) CreateTransaction(ctx context.Context, tx pgx.Tx, record dto.Transaction) (dto.Transaction, error) {
	queries := r.queries.WithTx(tx)
	transactionId, err := queries.CreateTransaction(ctx, ledgerdb.CreateTransactionParams{
		Identifier: record.Identifier,
		TransferID: record.TransferId,
		CreatedAt:  record.CreatedAt,
		Status:     record.Status,
	})
	if err != nil {
		return dto.Transaction{}, err
	}

	record.Id = transactionId
	return record, nil
}

func (r *LedgerPostgresRepository) GetAccountBalance(ctx context.Context, tx pgx.Tx, accountId int64) (dto.AccountBalance, error) {
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

func (r *LedgerPostgresRepository) UpdateAccountBalance(ctx context.Context, tx pgx.Tx, record dto.AccountBalance) error {
	queries := r.queries.WithTx(tx)
	balance, err := Float64ToNumeric(record.Availablebalance)
	if err != nil {
		return err
	}
	return queries.UpdateAccountBalance(ctx, ledgerdb.UpdateAccountBalanceParams{
		AccountID:        record.AccountId,
		AvailableBalance: balance,
		UpdatedAt:        record.UpdatedAt,
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
