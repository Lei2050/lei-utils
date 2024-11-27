package webreqproxy

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"

	wrproxy "github.com/Lei2050/lei-utils/webrequestproxy"
)

type HttpRequestRoutine struct {
	requestQ chan *wrproxy.HttpRequest
	close    chan bool
	proxy    *HttpRequestProxy
}

func (h *HttpRequestRoutine) Run() {
	for {
		select {
		case request := <-h.requestQ:
			h.SendHttpRequest(request)
		case <-h.close:
			for more := true; more; {
				select {
				case ptRequest := <-h.requestQ:
					h.SendHttpRequest(ptRequest)
				default:
					more = false
				}
			}
			return
		}
	}
}

func (h *HttpRequestRoutine) SendHttpRequest(request *wrproxy.HttpRequest) {
	httClient := http.Client{Timeout: 10 * time.Second}
	httpRsp := &wrproxy.HttpResponse{Resp: nil, Ret: wrproxy.ErrInvalidHttpMethod, Handler: request.Handler}
	if request.Method == wrproxy.HTTP_METHOD_GET {
		httpRsp.Resp, httpRsp.Ret = httClient.Get(request.Url)
	} else if request.Method == wrproxy.HTTP_METHOD_POST {
		httpRsp.Resp, httpRsp.Ret = httClient.PostForm(request.Url, request.Param)
	} else if request.Method == wrproxy.HTTP_METHOD_POSTJSON {
		content := request.Param.Get("json")
		httpRsp.Resp, httpRsp.Ret = httClient.Post(request.Url, "application/json", bytes.NewBuffer([]byte(content)))
	}
	h.proxy.C <- httpRsp
}

type HttpRequestProxy struct {
	requestGs []*HttpRequestRoutine
	index     int
	C         chan *wrproxy.HttpResponse
}

func NewHttpRequestProxy(routineNum int) *HttpRequestProxy {
	return &HttpRequestProxy{
		requestGs: make([]*HttpRequestRoutine, routineNum),
		index:     0,
		C:         make(chan *wrproxy.HttpResponse, 10240),
	}
}

func (h *HttpRequestProxy) Run() {
	for i := 0; i < len(h.requestGs); i++ {
		h.requestGs[i] = &HttpRequestRoutine{
			requestQ: make(chan *wrproxy.HttpRequest, 1024),
			close:    make(chan bool, 1),
			proxy:    h,
		}
		go h.requestGs[i].Run()
	}
}

func (h *HttpRequestProxy) SendHttpRequest(request *wrproxy.HttpRequest) {
	tmp := len(h.requestGs)
	if h.index < tmp {
		h.requestGs[h.index].requestQ <- request
		h.index++
		if h.index >= tmp {
			h.index = 0
		}
	}
}

func (h *HttpRequestProxy) SendRequest(
	method wrproxy.HttpMethodType,
	url string,
	param url.Values,
	handler wrproxy.HttpResponseHandler,
) {
	h.SendHttpRequest(&wrproxy.HttpRequest{
		Method:  method,
		Url:     url,
		Param:   param,
		Handler: handler,
	})
}

func (h *HttpRequestProxy) SendRequestWithCallback(
	method wrproxy.HttpMethodType,
	url string,
	param url.Values,
	cb wrproxy.HttpResponseCallback,
) {
	h.SendHttpRequest(&wrproxy.HttpRequest{
		Method:  method,
		Url:     url,
		Param:   param,
		Handler: cb,
	})
}

func (h *HttpRequestProxy) realProcResponse(response *wrproxy.HttpResponse, data []byte, err error) {
	switch handler := response.Handler.(type) {
	case wrproxy.HttpResponseCallback:
		handler(data, err)
	case wrproxy.HttpResponseHandler:
		handler.ProcHttpResponse(data, err)
	default:
		panic("HttpRequestProxy realProcResponse unknown type")
	}
}

func (h *HttpRequestProxy) ProcessResponse(response *wrproxy.HttpResponse) {
	if response == nil {
		return
	}
	if response.Ret != nil || response.Resp == nil {
		h.realProcResponse(response, nil, response.Ret)
		return
	}

	data, err := io.ReadAll(response.Resp.Body)
	defer response.Resp.Body.Close()
	if err != nil {
		response.Ret = err
		h.realProcResponse(response, nil, response.Ret)
		return
	}

	h.realProcResponse(response, data, nil)
}

func (h *HttpRequestProxy) Stop() {
	tmp := uint32(len(h.requestGs))
	for i := uint32(0); i < tmp; i++ {
		h.requestGs[i].close <- true
	}

	h.index = 0
	h.requestGs = make([]*HttpRequestRoutine, 0)
}
