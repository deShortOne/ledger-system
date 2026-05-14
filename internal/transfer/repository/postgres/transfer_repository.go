package postgres

import (
	"context"
	"errors"
	"strconv"

	"github.com/deshortone/ledger-system/internal/platform/database_base"
	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/deshortone/ledger-system/internal/transfer/repository/postgres/transferdb"
	"github.com/deshortone/ledger-system/pkg/failure"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransferPostgresRepository struct {
	database_base.BaseRepo
}

func NewTransferPostgresRepository(pool *pgxpool.Pool) TransferPostgresRepository {
	if pool == nil {
		panic("pool cannot be nil")
	}
	return TransferPostgresRepository{
		BaseRepo: database_base.NewBaseRepo(pool),
	}
}

func (r TransferPostgresRepository) CreateTransfer(ctx context.Context, request dto.NewTransfer) error {
	executor := r.GetExecutor(ctx)
	queries := transferdb.New(executor)

	return queries.CreateTransfer(ctx, transferdb.CreateTransferParams{
		Identifier:   request.Identifier,
		Identifier_2: request.TransferRequestId,
		ExecutedAt:   request.ExecutedAt.Time,
	})
}

func (r TransferPostgresRepository) CreateTransferRequest(ctx context.Context, request dto.NewTransferRequest) error {
	balance, err := Float64ToNumeric(request.Amount)
	if err != nil {
		return failure.NewFailure(
			failure.ConversionError,
			failure.GeneralFailure,
			err,
			"Failed to convert transfer request amount",
		)
	}

	executor := r.GetExecutor(ctx)
	queries := transferdb.New(executor)

	err = queries.CreateTransferRequest(ctx, transferdb.CreateTransferRequestParams{
		Identifier:   request.Identifier,
		Identifier_2: request.FromAccountId,
		Identifier_3: request.ToAccountId,
		Amount:       balance,
		Status:       request.Status,
		RequestedAt:  request.RequestedAt.Time,
	})
	if err != nil {
		return failure.NewFailure(
			failure.TransferRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to create transfer request",
		)
	}

	return nil
}

func (r TransferPostgresRepository) GetTransferStatus(ctx context.Context, id uuid.UUID) (string, string, error) {
	exector := r.GetExecutor(ctx)
	queries := transferdb.New(exector)

	data, err := queries.GetTranserRequestStatus(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", failure.NewFailure(
				failure.TransferRequestNotFound,
				failure.NotFound,
				err,
				"No transfer status record exists for the requested identifier",
			)
		}
		return "", "", failure.NewFailure(
			failure.TransferRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to fetch transfer status",
		)
	}

	return data.Status, data.FailureReason.String, nil
}

func (r TransferPostgresRepository) UpdateTransferRequestStatusWithTx(ctx context.Context, id uuid.UUID, status string) error {
	executor := r.GetExecutor(ctx)
	queries := transferdb.New(executor)

	err := queries.UpdateTransferRequestStatus(ctx, transferdb.UpdateTransferRequestStatusParams{
		Identifier:    id,
		Status:        status,
		FailureReason: pgtype.Text{Valid: false},
	})
	if err != nil {
		return failure.NewFailure(
			failure.TransferRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to update transfer request status",
		)
	}

	return nil
}

func (r TransferPostgresRepository) UpdateTransferRequestStatusWithFailure(ctx context.Context, id uuid.UUID, status, failureMessage string) error {
	executor := r.GetExecutor(ctx)
	queries := transferdb.New(executor)

	err := queries.UpdateTransferRequestStatus(ctx, transferdb.UpdateTransferRequestStatusParams{
		Identifier:    id,
		Status:        status,
		FailureReason: pgtype.Text{String: failureMessage, Valid: true},
	})
	if err != nil {
		return failure.NewFailure(
			failure.TransferRepositoryError,
			failure.GeneralFailure,
			err,
			"Failed to update transfer request failure status",
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
