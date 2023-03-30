package requester

type QueryParamOption struct {
	param string
}

func QueryParam(name, value string) *QueryParamOption {
	return &QueryParamOption{
		param: name + "=" + value,
	}
}

func (o *QueryParamOption) Apply(req *HTTPRequester) {
	if req.containsQueryParams {
		req.path += ("&" + o.param)
	} else {
		req.containsQueryParams = true
		req.path += ("?" + o.param)
	}
}
