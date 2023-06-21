package requester

type HeaderOption struct {
	key   string
	value string
}

func Header(key, value string) *HeaderOption {
	return &HeaderOption{
		key:   key,
		value: value,
	}
}

func (o *HeaderOption) Apply(req *HTTPRequester) {
	req.headers.Set(o.key, o.value)
}

func ContentType(value string) *HeaderOption {
	return Header("Content-Type", value)
}
