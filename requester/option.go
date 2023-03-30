package requester

type Option interface {
	Apply(r *HTTPRequester)
}
