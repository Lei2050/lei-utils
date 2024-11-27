package webreqproxy

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"

	wrproxy "github.com/Lei2050/lei-utils/webrequestproxy"
	wp "github.com/Lei2050/lei-utils/work_pool"
)

// 这里的HttpRequestProxy和raw_go目录下的HttpRequestProxy实现相同的功能。
// 都是http请求代理；后者是跑固定数量的携程。而这里采用的work_pool，可设定最大携程数量，
// 然后可以根据空闲与否，动态地增减携程数量，更为合理。
type HttpRequestProxy struct {
	wp *wp.WorkerPool
	C  chan *wrproxy.HttpResponse
}

func NewHttpRequestProxy(maxWorkers int) *HttpRequestProxy {
	return &HttpRequestProxy{
		wp: wp.NewWorkPool(maxWorkers),
		C:  make(chan *wrproxy.HttpResponse, 1024),
	}
}

func (h *HttpRequestProxy) SendHttpRequest(request *wrproxy.HttpRequest) {
	h.wp.Submit(func() {
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
		h.C <- httpRsp
	})
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
	h.wp.Stop()
}
