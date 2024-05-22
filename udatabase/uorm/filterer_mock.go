package uorm

import (
	"github.com/carlosarismendi/utils/udatabase/uorm/filters"
	"github.com/stretchr/testify/mock"
	"net/url"
)

type FiltererMock[T any] struct {
	mock.Mock
}

func (m *FiltererMock[T]) ParseFilters(values url.Values) ([]filters.ValuedFilter[T], error) {
	args := m.Called(values)
	res, _ := args.Get(0).([]filters.ValuedFilter[T])
	return res, args.Error(1)
}
