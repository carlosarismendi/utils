package filters

import (
	"github.com/carlosarismendi/utils/db/infrastructure/filters"
)

type TextFieldFilter struct {
	field string
}

func TextField(field string) *TextFieldFilter {
	return &TextFieldFilter{
		field: field,
	}
}

func (f *TextFieldFilter) Apply(values []string) (queryResult string, args []interface{}, rErr error) {
	return filters.ApplyTextField(f.field, values)
}
