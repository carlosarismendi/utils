package requester

import "net/http"

type DoerOption struct {
	doer *http.Client
}

func Doer(doer *http.Client) *DoerOption {
	return &DoerOption{
		doer: doer,
	}
}

func (o *DoerOption) Apply(req *HTTPRequester) {
	req.Doer = o.doer
}
