package transfer

import (
	"github.com/deshortone/ledger-system/internal/platform/database_base"
	"github.com/deshortone/ledger-system/internal/transfer/application"
	"github.com/deshortone/ledger-system/internal/transfer/controller"
	"github.com/deshortone/ledger-system/internal/transfer/domain"
	"github.com/deshortone/ledger-system/internal/transfer/repository/postgres"
	"github.com/deshortone/ledger-system/internal/transfer/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransferModule struct {
	Handler controller.Handler
}

func SetupTransferModule(
	pool *pgxpool.Pool,
	ledgerService domain.LedgerService,
) TransferModule {
	repository := postgres.NewTransferPostgresRepository(pool)
	service := service.NewTransferService(repository)
	uow := database_base.NewPgUnitOfWork(pool)
	transferApplication := application.NewTransferMoneyBetweenAccounts(service, ledgerService, uow)
	handler := controller.NewHandler(transferApplication, service)

	return TransferModule{
		Handler: handler,
	}
}
