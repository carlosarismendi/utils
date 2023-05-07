package requester

import (
	"encoding/json"
	"io"
	"net/http"
)

type HTTPRequester struct {
	url                 string
	path                string
	method              string
	headers             http.Header
	body                io.Reader
	containsQueryParams bool

	Doer *http.Client
}

func NewRequester(options ...Option) *HTTPRequester {
	r := &HTTPRequester{
		url:                 "",
		method:              "",
		headers:             make(http.Header),
		body:                nil,
		containsQueryParams: false,
		Doer:                http.DefaultClient,
	}

	return r.withOptions(options...)
}

func (r *HTTPRequester) Send(dst interface{}, options ...Option) (*http.Response, []byte, error) {
	if dstOpt, ok := dst.(Option); ok {
		options = append(options, nil)
		copy(options[1:], options)
		options[0] = dstOpt
		dst = nil
	}

	requesterClone := r.withOptions(options...)

	request, err := requesterClone.prepareRequest()
	if err != nil {
		return nil, nil, err
	}

	response, err := r.Doer.Do(request)
	if err != nil {
		return response, nil, err
	}

	body, err := requesterClone.readBody(dst, response)
	return response, body, err
}

func (r *HTTPRequester) prepareRequest() (*http.Request, error) {
	request, err := http.NewRequest(r.method, r.url+r.path, r.body)
	if err != nil {
		return nil, err
	}

	request.Header = r.headers

	return request, nil
}

func (r *HTTPRequester) readBody(dst interface{}, response *http.Response) ([]byte, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if dst != nil {
		err = json.Unmarshal(body, dst)
	}

	return body, err
}

func (r *HTTPRequester) withOptions(options ...Option) *HTTPRequester {
	clone := r.clone()
	for i := range options {
		opt := options[i]
		opt.Apply(clone)
	}
	return clone
}

func (r *HTTPRequester) clone() *HTTPRequester {
	clone := *r
	return &clone
}
