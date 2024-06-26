package udatabase

type ResourcePage[T any] struct {
	Total     int64 `json:"total"`
	Limit     int64 `json:"limit"`
	Offset    int64 `json:"offset"`
	Resources []T   `json:"resources"`
}
