package webreqproxy

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"

	wrproxy "github.com/Lei2050/lei-utils/webrequestproxy"
)

type HttpResponseHandlerSayHello struct {
	context string
}

func (h *HttpResponseHandlerSayHello) ProcHttpResponse(data []byte, err error) {
	fmt.Printf("SayHello cb context:%s, data:%s, err:%+v\n", h.context, string(data), err)
}

func TestProxy(t *testing.T) {
	proxy := NewHttpRequestProxy(20)
	proxy.Run()
	defer proxy.Stop()

	proxy.SendHttpRequest(&wrproxy.HttpRequest{Method: wrproxy.HTTP_METHOD_GET,
		Url:   "http://127.0.0.1:9090/sayhello?a=123&b=100",
		Param: nil,
		Handler: &HttpResponseHandlerSayHello{
			context: "This is testing !",
		},
	})

	proxy.SendRequest(wrproxy.HTTP_METHOD_GET,
		"http://127.0.0.1:9090/sayHi?a=666&b=100",
		nil,
		&HttpResponseHandlerSayHello{
			context: "This is sayHi !",
		},
	)

	params := url.Values{}
	params.Set("a", strconv.FormatInt(999, 10))
	params.Set("b", strconv.FormatInt(100, 10))
	proxy.SendRequestWithCallback(wrproxy.HTTP_METHOD_POST,
		"http://127.0.0.1:9090/saygoodbye",
		params,
		func(data []byte, err error) {
			fmt.Printf("Goodbye data:%s, err:%+v\n", string(data), err)
		},
	)

	for r := range proxy.C {
		proxy.ProcessResponse(r)
	}
}
