package requester

type ContentTypeOption struct {
	contentType string
}

func ContentType(contentType string) *ContentTypeOption {
	return &ContentTypeOption{
		contentType: contentType,
	}
}

func (o *ContentTypeOption) Apply(req *HTTPRequester) {
	req.contentType = o.contentType
}
