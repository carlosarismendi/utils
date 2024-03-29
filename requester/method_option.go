package requester

import "net/http"

type MethodOption struct {
	method string
	path   string
}

func newMethodOption(method, path string) *MethodOption {
	return &MethodOption{
		method: method,
		path:   path,
	}
}

func (o *MethodOption) Apply(req *HTTPRequester) {
	req.method = o.method
	req.path = o.path
}

func Get(path string) *MethodOption {
	return newMethodOption(http.MethodGet, path)
}

func Post(path string) *MethodOption {
	return newMethodOption(http.MethodPost, path)
}

func Put(path string) *MethodOption {
	return newMethodOption(http.MethodPut, path)
}

func Patch(path string) *MethodOption {
	return newMethodOption(http.MethodPatch, path)
}
