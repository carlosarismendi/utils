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
		req.url = req.url + "&" + o.param
	} else {
		req.containsQueryParams = true
		req.url = req.url + "?" + o.param
	}
}
