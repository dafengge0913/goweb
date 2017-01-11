package goweb

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/dafengge0913/golog"
)

type Context struct {
	log            *golog.Logger
	Req            *http.Request
	session        *session
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

func (ctx *Context) Param(key string) string {
	return ctx.params[key]
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
		ctx.log.Error("unmarshal json error: %v", err)
		return ctx.jsonParams
	}
	return ctx.jsonParams
}

func (ctx *Context) ResponseTemplate(tpl *template.Template, data interface{}) {
	if err := tpl.Execute(ctx.responseWriter, data); err != nil {
		ctx.log.Error("template execute error: %v", err)
	}
}

func (ctx *Context) ResponseJSON(data interface{}) {
	if jsonBytes, err := json.Marshal(data); err != nil {
		ctx.log.Error("marshal json error: %v", err)
	} else {
		ctx.responseWriter.Write(jsonBytes)
	}
}

func (ctx *Context) SetCookie(name, value string) {
	c := &http.Cookie{
		Name:  name,
		Value: value,
	}
	http.SetCookie(ctx.responseWriter, c)
}

func (ctx *Context) SetRawCookie(c *http.Cookie) {
	http.SetCookie(ctx.responseWriter, c)
}

func (ctx *Context) Cookie(name string) *http.Cookie {
	c, err := ctx.Req.Cookie(name)
	if err != nil {
		return nil
	}
	return c
}

func (ctx *Context) DelCookie(name string) {
	c := ctx.Cookie(name)
	if c != nil {
		c.MaxAge = -1
		http.SetCookie(ctx.responseWriter, c)
	}
}

func (ctx *Context) Session() *session {
	if ctx.session != nil {
		return ctx.session
	}
	ctx.log.Warn("session is not enable")
	return nil
}
