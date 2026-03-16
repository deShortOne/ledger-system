package service

import (
	"context"

	"github.com/deshortone/ledger-system/internal/ledger/contracts"
	"github.com/deshortone/ledger-system/internal/ledger/domain"
	"github.com/deshortone/ledger-system/internal/ledger/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LedgerService struct {
	ledgerRepository         domain.LedgerRepository
	accountBalanceRepository domain.AccountBalanceRepository
}

func NewLedgerService(
	ledgerRepository domain.LedgerRepository,
	accountBalanceRepository domain.AccountBalanceRepository,
) LedgerService {
	return LedgerService{
		ledgerRepository:         ledgerRepository,
		accountBalanceRepository: accountBalanceRepository,
	}
}

func (s LedgerService) AddToLedger(ctx context.Context, tx pgx.Tx, request contracts.AddToLedgerRequest) error {
	var sumOfMonies float64
	transactionId := uuid.New()
	err := s.ledgerRepository.CreateTransaction(ctx, tx, dto.Transaction{
		Identifier: uuid.New(),
		TransferId: request.TransferId,
		CreatedAt:  request.CreatedAt,
		Status:     "posted",
	})
	if err != nil {
		return err
	}

	for _, entry := range request.Entries {
		accountBalance, err := s.accountBalanceRepository.GetAccountBalance(ctx, tx, entry.AccountId)
		if err != nil {
			return err
		}

		if entry.Direction == contracts.CREDIT {
			sumOfMonies += entry.Amount
			accountBalance.Availablebalance += entry.Amount
		} else {
			if entry.Amount > accountBalance.Availablebalance {
				return contracts.ErrOneOfTheAccountsDoNotHaveEnoughMoney
			}

			sumOfMonies -= entry.Amount
			accountBalance.Availablebalance -= entry.Amount
		}
		accountBalance.UpdatedAt = request.CreatedAt

		err = s.ledgerRepository.CreateLedgerEntry(ctx, tx, dto.LedgerEntry{
			Identifier:    uuid.New(),
			TransactionId: transactionId,
			AccountId:     entry.AccountId,
			Amount:        entry.Amount,
			Direction:     contracts.LedgerDirection(entry.Direction),
			CreatedAt:     request.CreatedAt,
		})
		if err != nil {
			return err
		}

		if err = s.accountBalanceRepository.UpdateAccountBalance(ctx, tx, accountBalance); err != nil {
			return err
		}
	}
	if sumOfMonies != 0 {
		return contracts.ErrDoubleEntryViolated
	}

	return nil
}
