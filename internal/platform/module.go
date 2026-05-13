package platform

import (
	"github.com/deshortone/ledger-system/internal/platform/controller"
	"github.com/deshortone/ledger-system/internal/platform/repository"
	"github.com/deshortone/ledger-system/internal/platform/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlatformModule struct {
	PlatformHandler controller.Handler
}

func SetupPlatformModule(pool *pgxpool.Pool) PlatformModule {
	platformRepository := repository.NewPlatformPostgresRepository(pool)
	platformService := service.NewPlatformService(platformRepository)

	return PlatformModule{
		PlatformHandler: controller.NewHandler(platformService),
	}
}
