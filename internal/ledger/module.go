package ledger

import (
	"github.com/deshortone/ledger-system/internal/ledger/repository/postgres"
	"github.com/deshortone/ledger-system/internal/ledger/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerModule struct {
	LedgerService service.LedgerService
}

func SetupLedgerModule(pool *pgxpool.Pool) LedgerModule {
	ledgerRepository := postgres.NewLedgerPostgresRepository(pool)
	accountBalanceRepository := postgres.NewAccountBalancePostgresRepository(pool)

	return LedgerModule{
		LedgerService: service.NewLedgerService(ledgerRepository, accountBalanceRepository),
	}
}
