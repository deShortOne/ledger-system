package domain

import "context"

type PlatformHealth interface {
	IsUp(ctx context.Context) error
}
