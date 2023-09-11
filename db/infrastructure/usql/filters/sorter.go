package filters

import (
	"github.com/carlosarismendi/utils/db/infrastructure/filters"
)

type SortFilter struct {
	allowedFields map[string]bool
}

func Sort(fields ...string) *SortFilter {
	allowedFields := make(map[string]bool)
	for _, f := range fields {
		allowedFields[f] = true
	}
	return &SortFilter{
		allowedFields: allowedFields,
	}
}

func (f *SortFilter) Apply(values []string) (string, error) {
	return filters.ApplySorter(f.allowedFields, values)
}
