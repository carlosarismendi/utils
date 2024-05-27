package requester

import "strings"

type AppendPathOption struct {
	path string
}

func AppendPath(path string) *AppendPathOption {
	path = strings.TrimLeft(path, "/")
	return &AppendPathOption{
		path: path,
	}
}

func (o *AppendPathOption) Apply(req *HTTPRequester) {
	if req.url != "" && req.url[len(req.url)-1] != '/' &&
		(req.path == "" || req.path[len(req.path)-1] != '/') {
		req.path += "/"
	}
	req.path += o.path
}
