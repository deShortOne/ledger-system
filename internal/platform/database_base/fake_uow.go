package database_base

import "context"

type FakeUOW struct{}

func (f FakeUOW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
