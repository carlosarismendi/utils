package requester

type AppendPathOption struct {
	path string
}

func AppendPath(path string) *AppendPathOption {
	return &AppendPathOption{
		path: path,
	}
}

func (o *AppendPathOption) Apply(req *HTTPRequester) {
	req.path += o.path
}
