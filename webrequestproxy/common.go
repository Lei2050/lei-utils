package webrequestproxy

import "errors"

type HttpMethodType = int

const (
	HttpMethodGet HttpMethodType = iota
	HttpMethodPost
	HttpMethodPostJson
)

var (
	ErrInvalidHttpMethod error = errors.New("invalid http method")
)
