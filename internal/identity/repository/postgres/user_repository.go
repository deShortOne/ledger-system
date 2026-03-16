package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/deshortone/ledger-system/internal/identity/common"
	"github.com/deshortone/ledger-system/internal/identity/dto"
	userdb "github.com/deshortone/ledger-system/internal/identity/repository/postgres/user_db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgresRepository struct {
	queries *userdb.Queries
}

func NewUserPostgresRepository(pool *pgxpool.Pool) *UserPostgresRepository {
	if pool == nil {
		panic("pool cannot be nil")
	}

	return &UserPostgresRepository{
		queries: userdb.New(pool),
	}
}

func (r *UserPostgresRepository) CreateUser(ctx context.Context, user dto.User) (dto.User, error) {
	userId, err := r.queries.CreateUser(ctx, userdb.CreateUserParams{
		Identifier: user.Identifier,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
	})
	if err != nil {
		return dto.User{}, err
	}

	user.Id = userId

	return user, nil
}

func (r *UserPostgresRepository) GetUser(ctx context.Context, identifier uuid.UUID) (dto.User, error) {
	user, err := r.queries.GetUser(ctx, identifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.User{}, common.ErrUserIdentifierNotFound
		}
		return dto.User{}, err
	}

	return dto.User{
		Id:         user.ID,
		Identifier: identifier,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
	}, nil
}
