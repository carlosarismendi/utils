package requester

type URLOption struct {
	url string
}

func URL(url string) *URLOption {
	return &URLOption{
		url: url,
	}
}

func (o *URLOption) Apply(req *HTTPRequester) {
	req.url = o.url
}
