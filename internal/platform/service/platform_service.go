package service

import (
	"context"

	"github.com/deshortone/ledger-system/internal/platform/domain"
)

type PlatformService struct {
	repository domain.PlatformHealth
}

func NewPlatformService(repository domain.PlatformHealth) PlatformService {
	return PlatformService{
		repository: repository,
	}
}

func (s PlatformService) IsUp(ctx context.Context) error {
	return s.repository.IsUp(ctx)
}
