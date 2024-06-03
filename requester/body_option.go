package requester

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/carlosarismendi/utils/uerr"
)

type BodyReaderOption struct {
	body io.Reader
}

func BodyReader(body io.Reader) *BodyReaderOption {
	return &BodyReaderOption{
		body: body,
	}
}

func (o *BodyReaderOption) Apply(req *HTTPRequester) {
	req.body = o.body
}

func Body(body any) *BodyReaderOption {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		pErr := uerr.NewError(uerr.GenericError,
			"Error marshaling body to JSON in http request").WithCause(err)
		panic(pErr)
	}

	return BodyReader(&buf)
}
