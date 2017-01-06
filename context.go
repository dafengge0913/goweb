package goweb

import "net/http"

type Context struct {
	Req    *http.Request
	Params map[string]string
}

func newContext(req *http.Request) *Context {
	return &Context{
		Req:    req,
		Params: make(map[string]string),
	}
}
