package filters

type Sorter interface {
	Apply(values []string) (string, error)
}

type Filter interface {
	Apply(values []string) (string, []interface{}, error)
}
