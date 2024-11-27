package webrequestproxy

import "errors"

type HttpMethodType = int

const (
	HTTP_METHOD_GET HttpMethodType = iota
	HTTP_METHOD_POST
	HTTP_METHOD_POSTJSON
)

var (
	ErrInvalidHttpMethod error = errors.New("invalid http method")
)
