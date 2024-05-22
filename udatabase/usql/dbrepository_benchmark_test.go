package usql

import (
	"net/url"
	"testing"

	usqlFilters "github.com/carlosarismendi/utils/udatabase/usql/filters"
)

type filtersBenchmark struct {
	name    string
	filters url.Values
}

func (fb *filtersBenchmark) getFilters() url.Values {
	cv := url.Values{}
	for k, v := range fb.filters {
		cv[k] = v
	}

	return cv
}

func BenchmarkApplyFilters(b *testing.B) {
	var table = []filtersBenchmark{
		{
			name: "1Text",
			filters: createFilters(
				newFilter("name", "car"),
			),
		},
		{
			name: "1Number",
			filters: createFilters(
				newFilter("random_number", "1"),
			),
		},
		{
			name: "limit",
			filters: createFilters(
				newFilter("limit", "2"),
			),
		},
		{
			name: "offset",
			filters: createFilters(
				newFilter("offset", "2"),
			),
		},
		{
			name: "LimitOffset2Number1Text",
			filters: createFilters(
				newFilter("limit", "10"),
				newFilter("offset", "5"),
				newFilter("random_number", "1", "2"),
				newFilter("name", "car"),
			),
		},
		{
			name: "LimitOffset1Number2Text",
			filters: createFilters(
				newFilter("limit", "10"),
				newFilter("offset", "5"),
				newFilter("random_number", "2"),
				newFilter("name", "car", "airplane"),
			),
		},
		{
			name: "SortTextAscNumberDesc",
			filters: createFilters(
				newFilter("sort", "name", "-random_number"),
			),
		},
	}

	var dbHolder = NewTestDBHolder("db_usql_repository_test_get_context")
	var filtersMap = map[string]usqlFilters.Filter{
		"id":            usqlFilters.TextField("id"),
		"name":          usqlFilters.TextField("name"),
		"random_number": usqlFilters.NumField("random_number"),
	}

	var sortersMap = map[string]usqlFilters.Sorter{
		"sort": usqlFilters.Sort("name", "random_number"),
	}

	var r = NewDBRepository[*Resource](dbHolder.DBHolder, filtersMap, sortersMap)
	var db = r.GetDBInstance()

	var v url.Values
	for _, fb := range table {
		b.Run(fb.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				v = fb.getFilters()
				_, _, _, _, _ = r.ApplyFilters(db, "SELECT id, name, random_number AS RandomNumber FROM resources", v)
			}
		})
	}
}
