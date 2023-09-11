package filters

import (
	"github.com/carlosarismendi/utils/db/infrastructure/filters"
)

type NumFieldFilter struct {
	field string
}

func NumField(field string) *NumFieldFilter {
	return &NumFieldFilter{
		field: field,
	}
}

func (f *NumFieldFilter) Apply(values []string) (queryResult string, args []interface{}, rErr error) {
	return filters.ApplyNumField(f.field, values)
}
