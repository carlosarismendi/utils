package requester

import "net/url"

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

type QueryParamsOption struct {
	params []*QueryParamOption
}

func QueryParams(v url.Values) *QueryParamsOption {
	params := make([]*QueryParamOption, 0, len(v))
	for k, v := range v {
		params = append(params, QueryParam(k, v[0]))
	}

	return &QueryParamsOption{
		params: params,
	}
}

func (o *QueryParamsOption) Apply(req *HTTPRequester) {
	for _, p := range o.params {
		p.Apply(req)
	}
}
