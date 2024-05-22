package udatabase

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type TransactionalMock struct {
	mock.Mock
}

func (m *TransactionalMock) Begin(ctx context.Context) (context.Context, error) {
	args := m.Called(ctx)
	res, _ := args.Get(0).(context.Context)
	return res, args.Error(1)
}

func (m *TransactionalMock) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *TransactionalMock) Rollback(ctx context.Context) {
	_ = m.Called(ctx)
}
