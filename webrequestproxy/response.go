package webrequestproxy

import "net/http"

type HttpResponseCallback func([]byte, error)

type HttpResponseHandler interface {
	ProcHttpResponse([]byte, error)
}

type HttpResponse struct {
	Resp *http.Response
	Ret  error
	//handler HttpResponseHandler
	Handler any
}
