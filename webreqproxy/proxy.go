package webreqproxy

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"

	wp "github.com/Lei2050/lei-utils/work_pool"
)

// 这里的HttpRequestProxy和webrequestproxy目录下的HttpRequestProxy实现相同的功能。
// 都是http请求代理；后者是跑固定数量的携程。而这里采用的work_pool，可设定最大携程数量，
// 然后可以根据空闲与否，动态地增减携程数量，更为合理。
type HttpRequestProxy struct {
	wp *wp.WorkerPool
	C  chan *HttpResponse
}

func NewHttpRequestProxy(maxWorkers int) *HttpRequestProxy {
	return &HttpRequestProxy{
		wp: wp.NewWorkPool(maxWorkers),
		C:  make(chan *HttpResponse, 1024),
	}
}

func (h *HttpRequestProxy) SendHttpRequest(request *HttpRequest) {
	h.wp.Submit(func() {
		httClient := http.Client{Timeout: 10 * time.Second}
		httpRsp := &HttpResponse{nil, ErrInvalidHttpMethod, request.handler}
		if request.Method == HTTP_METHOD_GET {
			httpRsp.Resp, httpRsp.Ret = httClient.Get(request.Url)
		} else if request.Method == HTTP_METHOD_POST {
			httpRsp.Resp, httpRsp.Ret = httClient.PostForm(request.Url, request.Param)
		} else if request.Method == HTTP_METHOD_POSTJSON {
			content := request.Param.Get("json")
			httpRsp.Resp, httpRsp.Ret = httClient.Post(request.Url, "application/json", bytes.NewBuffer([]byte(content)))
		}
		h.C <- httpRsp
	})
}

func (h *HttpRequestProxy) SendRequest(
	method HttpMethodType,
	url string,
	param url.Values,
	handler HttpResponseHandler,
) {
	h.SendHttpRequest(&HttpRequest{
		Method:  method,
		Url:     url,
		Param:   param,
		handler: handler,
	})
}

func (h *HttpRequestProxy) SendRequestWithCallback(
	method HttpMethodType,
	url string,
	param url.Values,
	cb HttpResponseCallback,
) {
	h.SendHttpRequest(&HttpRequest{
		Method:  method,
		Url:     url,
		Param:   param,
		handler: cb,
	})
}

func (h *HttpRequestProxy) realProcResponse(response *HttpResponse, data []byte, err error) {
	switch handler := response.handler.(type) {
	case HttpResponseCallback:
		handler(data, err)
	case HttpResponseHandler:
		handler.ProcHttpResponse(data, err)
	default:
		panic("HttpRequestProxy realProcResponse unknown type")
	}
}

func (h *HttpRequestProxy) ProcessResponse(response *HttpResponse) {
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
