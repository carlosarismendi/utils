package domain

type ResourcePage struct {
	Total int64 `json:"total"`
	Limit int64 `json:"limit"`
	// TODO: implement offset in DBRepositoryFind
	Offset int64 `json:"offset"`

	// Resource will be a pointer to the type pased as
	// dst parameter in Find method. In this example,
	// *[]*Resource.
	Resources interface{} `json:"resources"`
}
