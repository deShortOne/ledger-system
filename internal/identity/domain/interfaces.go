package domain

import (
	"context"

	"github.com/deshortone/ledger-system/internal/identity/dto"
	"github.com/google/uuid"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account dto.Account) (dto.Account, error)
	GetAccountsOwnedByUser(ctx context.Context, user dto.User) ([]dto.Account, error)
}

type AccountService interface {
	AddAccountToUser(ctx context.Context, userIdentifier uuid.UUID, accountToCreate dto.Account) (dto.Account, error)
	GetAccountsOwnedByUser(ctx context.Context, identifier uuid.UUID) ([]dto.Account, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user dto.User) (dto.User, error)
	GetUser(ctx context.Context, identifier uuid.UUID) (dto.User, error)
}

type UserService interface {
	CreateNewUser(ctx context.Context, firstName, lastName string) (dto.User, error)
}
