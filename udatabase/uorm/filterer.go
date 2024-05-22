package uorm

import (
	"github.com/carlosarismendi/utils/udatabase/uorm/filters"
	"net/url"
)

type Filterer[T any] interface {
	ParseFilters(values url.Values) ([]filters.ValuedFilter[T], error)
}
