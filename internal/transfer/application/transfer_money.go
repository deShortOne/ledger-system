package application

import (
	"context"

	"github.com/deshortone/ledger-system/internal/ledger/contracts"
	"github.com/deshortone/ledger-system/internal/platform/database_base"
	"github.com/deshortone/ledger-system/internal/transfer/domain"
	"github.com/deshortone/ledger-system/internal/transfer/dto"
	"github.com/google/uuid"
)

type TransferMoneyBetweenAccounts struct {
	transferService domain.TransferService
	ledgerService   domain.LedgerService
	uow             database_base.UnitOfWork
}

func NewTransferMoneyBetweenAccounts(
	transferService domain.TransferService,
	ledgerService domain.LedgerService,
	uow database_base.UnitOfWork,
) *TransferMoneyBetweenAccounts {
	return &TransferMoneyBetweenAccounts{
		transferService: transferService,
		ledgerService:   ledgerService,
		uow:             uow,
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
	err = a.uow.Do(ctx, func(ctx1 context.Context) error {
		transferId, err := a.transferService.CreateTransfer(ctx1, transferRequestId, dto.NewCustomTimeNow())
		if err != nil {
			return err
		}

		err = a.ledgerService.AddToLedger(ctx1, contracts.AddToLedgerRequest{
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
		if err != nil {
			return err
		}

		return a.transferService.UpdateTransferRequestStatus(ctx1, transferRequestId, "Success")
	})
	if err != nil {
		return a.transferService.UpdateTransferRequestStatusWithFailure(ctx, transferRequestId, "Failed", "For reasons")
	}

	return nil
}
