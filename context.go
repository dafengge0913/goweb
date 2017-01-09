package goweb

import (
	"net/http"
	"html/template"
	"github.com/dafengge0913/golog"
)

type Context struct {
	log            *golog.Logger
	Req            *http.Request
	Params         map[string]string
	responseWriter http.ResponseWriter
}

func newContext(log *golog.Logger, req *http.Request, rw http.ResponseWriter) *Context {
	return &Context{
		log:           log,
		Req:           req,
		Params:        make(map[string]string),
		responseWriter:rw,
	}
}

func (ctx *Context) ResponseTemplate(tpl *template.Template, data interface{}) {
	if err := tpl.Execute(ctx.responseWriter, data); err != nil {
		ctx.log.Error("tpl execute error: %v", err)
	}
}
