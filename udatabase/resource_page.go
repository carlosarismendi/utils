package udatabase

type ResourcePage struct {
	Total     int64       `json:"total"`
	Limit     int64       `json:"limit"`
	Offset    int64       `json:"offset"`
	Resources interface{} `json:"resources"`
}
