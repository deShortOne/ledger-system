package application

import (
	"context"

	"github.com/deshortone/ledger-system/internal/ledger/contracts"
	"github.com/deshortone/ledger-system/internal/transfer/domain"
	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransferMoneyBetweenAccounts struct {
	transferService domain.TransferService
	pool            *pgxpool.Pool
	ledgerService   domain.LedgerService
}

func NewTransferMoneyBetweenAccounts(
	transferService domain.TransferService,
	pool *pgxpool.Pool,
	ledgerService domain.LedgerService,
) *TransferMoneyBetweenAccounts {
	return &TransferMoneyBetweenAccounts{
		transferService: transferService,
		pool:            pool,
		ledgerService:   ledgerService,
	}
}

// !! DO NOT USE YET - NOT ONLY UNTESTED BUT WRONG! SKELETON HERE ONLY
func (a *TransferMoneyBetweenAccounts) TransferMoney(ctx context.Context, fromAccountId, toAccountId uuid.UUID, amount float64) error {
	transferRequestId, err := a.transferService.CreateTransferRequest(ctx, dto.CreateNewTransfer{
		FromAccountId: fromAccountId,
		ToAccountId:   toAccountId,
		Amount:        amount,
		RequestedAt:   dto.NewCustomTimeNow(),
	})
	if err != nil {
		return err
	}

	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback(ctx)
		a.transferService.UpdateTransferRequestStatusWithFailure(ctx, transferRequestId, "Failed", "due to rollback")
	}()

	transferId, err := a.transferService.CreateTransfer(ctx, tx, transferRequestId, dto.NewCustomTimeNow())

	a.ledgerService.AddToLedger(ctx, tx, contracts.AddToLedgerRequest{
		TransferId: transferId,
		CreatedAt:  dto.NewCustomTimeNow().Time,
		Entries: []contracts.LedgerEntries{
			{
				AccountId: fromAccountId,
				Amount:    amount,
				Direction: contracts.DEBIT,
			},
			{
				AccountId: toAccountId,
				Amount:    amount,
				Direction: contracts.CREDIT,
			},
		},
	})

	a.transferService.UpdateTransferRequestStatus(ctx, tx, transferRequestId, "Success")

	if err = tx.Commit(ctx); err != nil {
		a.transferService.UpdateTransferRequestStatusWithFailure(ctx, transferRequestId, "Failed", "due to unable to commit")
		return err
	}

	return nil
}
