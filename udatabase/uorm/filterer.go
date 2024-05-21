package uorm

import (
	"github.com/carlosarismendi/utils/udatabase/uorm/filters"
	"net/url"
)

type Filterer interface {
	ParseFilters(values url.Values) ([]filters.ValuedFilter, error)
}
