package identity

import (
	"github.com/deshortone/ledger-system/internal/identity/controller"
	"github.com/deshortone/ledger-system/internal/identity/repository/postgres"
	"github.com/deshortone/ledger-system/internal/identity/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IdentityModule struct {
	Handler        controller.Handler
	AccountService service.AccountService
	UserService    service.UserService
}

func SetupIdentityModule(pool *pgxpool.Pool) IdentityModule {
	accountRepository := postgres.NewAccountPostgresRepository(pool)
	userRepository := postgres.NewUserPostgresRepository(pool)

	accountService := service.NewAccountService(accountRepository, userRepository)
	userService := service.NewUserService(userRepository)

	return IdentityModule{
		Handler:        controller.NewHandler(accountService, userService),
		AccountService: accountService,
		UserService:    userService,
	}
}
