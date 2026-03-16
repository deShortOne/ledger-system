package postgres

import (
	"context"
	"strconv"

	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/deshortone/ledger-system/internal/transfer/repository/postgres/transferdb"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransferPostgresRepository struct {
	queries *transferdb.Queries
}

func NewTransferPostgresRepository(pool *pgxpool.Pool) TransferPostgresRepository {
	if pool == nil {
		panic("pool cannot be nil")
	}
	return TransferPostgresRepository{
		queries: transferdb.New(pool),
	}
}

func (r TransferPostgresRepository) CreateTransfer(ctx context.Context, tx pgx.Tx, request dto.NewTransfer) error {
	queries := r.queries.WithTx(tx)

	return queries.CreateTransfer(ctx, transferdb.CreateTransferParams{
		Identifier:   request.Identifier,
		Identifier_2: request.TransferRequestId,
		ExecutedAt:   request.ExecutedAt.Time,
	})
}

func (r TransferPostgresRepository) CreateTransferRequest(ctx context.Context, request dto.NewTransferRequest) error {
	balance, err := Float64ToNumeric(request.Amount)
	if err != nil {
		return err
	}

	return r.queries.CreateTransferRequest(ctx, transferdb.CreateTransferRequestParams{
		Identifier:   request.Identifier,
		Identifier_2: request.FromAccountId,
		Identifier_3: request.ToAccountId,
		Amount:       balance,
		Status:       request.Status,
		RequestedAt:  request.RequestedAt.Time,
	})
}

func (r TransferPostgresRepository) UpdateTransferRequestStatusWithTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, status string) error {
	queries := r.queries.WithTx(tx)
	return queries.UpdateTransferRequestStatus(ctx, transferdb.UpdateTransferRequestStatusParams{
		Identifier:    id,
		Status:        status,
		FailureReason: pgtype.Text{Valid: false},
	})
}

func (r TransferPostgresRepository) UpdateTransferRequestStatusWithFailure(ctx context.Context, id uuid.UUID, status, failure string) error {
	return r.queries.UpdateTransferRequestStatus(ctx, transferdb.UpdateTransferRequestStatusParams{
		Identifier:    id,
		Status:        status,
		FailureReason: pgtype.Text{String: failure, Valid: true},
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
