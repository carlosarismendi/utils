package udatabase

import (
	"context"
)

type Transactional interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
}

func BeginTx(ctx context.Context, r Transactional) (context.Context, error) {
	return r.Begin(ctx)
}

// nolint:gocritic // rErr *error is required to be pointer to capture properly
// the errors returned by functions
func EndTx(ctx context.Context, r Transactional, rErr *error) {
	pErr := recover()
	if pErr != nil {
		r.Rollback(ctx)
		panic(pErr)
	}

	if rErr != nil && *rErr != nil {
		r.Rollback(ctx)
		return
	}

	err := r.Commit(ctx)
	if err != nil {
		r.Rollback(ctx)
	}
}
