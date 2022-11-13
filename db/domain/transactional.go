package domain

import "context"

type Transactional interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func BeginTx(ctx context.Context, r Transactional) (context.Context, error) {
	return r.Begin(ctx)
}

func EndTx(ctx context.Context, r Transactional, rErr *error) {
	if rErr != nil && *rErr != nil {
		_ = r.Rollback(ctx)
		return
	}

	defer func() {
		rErr := recover()
		if rErr != nil {
			_ = r.Rollback(ctx)
			panic(rErr)
		}

		_ = r.Commit(ctx)
	}()
}
