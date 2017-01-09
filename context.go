package goweb

import (
	"encoding/json"
	"github.com/dafengge0913/golog"
	"html/template"
	"net/http"
)

type Context struct {
	log            *golog.Logger
	Req            *http.Request
	params         map[string]string
	responseWriter http.ResponseWriter
	jsonParams     map[string]string
}

func newContext(log *golog.Logger, req *http.Request, rw http.ResponseWriter) *Context {
	return &Context{
		log:            log,
		Req:            req,
		params:         make(map[string]string),
		responseWriter: rw,
	}
}

func (ctx *Context) Params() map[string]string {
	return ctx.params
}

func (ctx *Context) JSONParams() map[string]string {
	if ctx.jsonParams != nil {
		return ctx.jsonParams
	}
	ctx.jsonParams = make(map[string]string)
	var jsonStr string
	for k := range ctx.Params() {
		jsonStr = k
		break
	}

	if len(jsonStr) == 0 {
		return ctx.jsonParams
	}

	if err := json.Unmarshal([]byte(jsonStr), &ctx.jsonParams); err != nil {
		ctx.log.Error("Unmarshal json error: %v", err)
		return ctx.jsonParams
	}
	return ctx.jsonParams
}

func (ctx *Context) ResponseTemplate(tpl *template.Template, data interface{}) {
	if err := tpl.Execute(ctx.responseWriter, data); err != nil {
		ctx.log.Error("tpl execute error: %v", err)
	}
}

func (ctx *Context) ResponseJSON(data interface{}) {
	if jsonBytes, err := json.Marshal(data); err != nil {
		ctx.log.Error("marshal json error: %v", err)
	} else {
		ctx.responseWriter.Write(jsonBytes)
	}
}
