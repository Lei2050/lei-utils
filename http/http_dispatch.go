package http

import (
	"fmt"
	"net/http"
	"net/url"
)

type HttpHandler interface {
	Execute(http.ResponseWriter, url.Values) bool
}

type HttpHandlerFactory interface {
	CreateHttpHandler() HttpHandler
}

type HttpHandlerMgr struct {
	cmdm map[string]HttpHandlerFactory
}

func NewHttpHandlerMgr() *HttpHandlerMgr {
	return &HttpHandlerMgr{
		cmdm: make(map[string]HttpHandlerFactory),
	}
}

func (hhm *HttpHandlerMgr) Register(name string, cmd HttpHandlerFactory) {
	hhm.cmdm[name] = cmd
}

func (hhm *HttpHandlerMgr) Dispatcher(action string, w http.ResponseWriter, r url.Values) bool {
	if cmd, exist := hhm.cmdm[action]; exist {
		h := cmd.CreateHttpHandler()
		if h == nil {
			panic(fmt.Errorf("invalid handler created by handler factory of %s", action))
		}
		h.Execute(w, r)
		return true
	}
	//HttpOutput(w, RET_PLATFORM_ACTION_NOT_FIND, "action not found:"+action, nil)
	w.Write([]byte("no find action"))
	return false
}
