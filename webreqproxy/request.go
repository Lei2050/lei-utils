package webreqproxy

import "net/url"

type HttpRequest struct {
	Method HttpMethodType
	Url    string
	Param  url.Values
	//handler HttpResponseHandler
	handler interface{}
}
