package requester

import (
	"bytes"
	"encoding/json"

	"github.com/carlosarismendi/utils/uerr"
)

type BodyOption struct {
	body interface{}
}

func Body(body interface{}) *BodyOption {
	return &BodyOption{
		body: body,
	}
}

func (o *BodyOption) Apply(req *HTTPRequester) {
	if o.body == nil {
		return
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(o.body)
	if err != nil {
		pErr := uerr.NewError(uerr.GenericError,
			"Error marshaling body to JSON in http request").WithCause(err)
		panic(pErr)
	}

	req.body = &buf
}
