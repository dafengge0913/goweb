package goweb

import (
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/dafengge0913/golog"
)

//todo provide set func
var defaultSessionManagerConfig = &sessionManagerConfig{
	sessionMaxAge: 600, // 10 minute
	sessionKey:    "sid",
	sessionIdLen:  16,
	cleanInterval: time.Second * 10,
}

type HTTPHandler func(ctx *Context)

type Server struct {
	ln      net.Listener
	log     *golog.Logger
	routers []iRouter
	sm      *sessionManager
}

func NewServer() *Server {
	return &Server{
		log:     golog.NewLogger(golog.LEVEL_DEBUG, nil),
		routers: make([]iRouter, 0),
		sm:      newSessionManager(defaultSessionManagerConfig),
	}
}

type iRouter interface {
	match(path string) bool
	matchIndex(path string) int
	handle(ctx *Context)
}

type baseRouter struct {
	re *regexp.Regexp
}

func (br *baseRouter) match(path string) bool {
	return br.re.MatchString(path)
}

func (br *baseRouter) matchIndex(path string) int {
	match := br.re.FindStringSubmatchIndex(path)
	matchLen := len(match)
	if matchLen == 0 {
		return -1
	}
	return match[matchLen-1]
}

type router struct {
	*baseRouter
	handler HTTPHandler
}

func (r *router) handle(ctx *Context) {
	r.handler(ctx)
}

type staticRouter struct {
	*baseRouter
	handler http.Handler
}

func (sr *staticRouter) handle(ctx *Context) {
	sr.handler.ServeHTTP(ctx.responseWriter, ctx.Req)
}

// start http Server
// call after all routers have been added
func (srv *Server) Serve(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv.ln = ln
	defer srv.Close()
	handler := http.NewServeMux()
	handler.Handle("/", srv)
	http.Serve(ln, handler)
	return nil
}

func (srv *Server) Close() {
	if srv.ln != nil {
		if err := srv.ln.Close(); err != nil {
			srv.log.Info("close server error: %v", err)
		}
	}
	if srv.sm != nil {
		srv.sm.Close()
	}
}

func (srv *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	if err := req.ParseForm(); err != nil {
		srv.log.Error("parse request form error: %v", err)
		return
	}
	ctx := newContext(srv.log, req, rw)
	if srv.sm != nil {
		srv.setSession(ctx)
	}
	req.ParseForm()
	if len(req.Form) > 0 {
		for k, v := range req.Form {
			ctx.params[k] = v[0]
		}
	}
	r := srv.findRouter(reqPath)
	if r == nil {
		srv.log.Debug("router not found path:%s", reqPath)
		return
	}
	r.handle(ctx)
}

func (srv *Server) findRouter(path string) iRouter {
	var mostMatch iRouter
	mostMatchIdx := 0
	for _, r := range srv.routers {
		if !r.match(path) {
			continue
		}
		matchIdx := r.matchIndex(path)
		if matchIdx >= mostMatchIdx { // use the last setting
			mostMatchIdx = matchIdx
			mostMatch = r
		}

	}
	return mostMatch
}

func (srv *Server) AddRouter(match string, handler HTTPHandler) {
	re, err := regexp.Compile(match)
	if err != nil {
		srv.log.Error("compile regexp error: %v", err)
		return
	}
	r := &router{
		baseRouter: &baseRouter{
			re: re,
		},
		handler: handler,
	}
	srv.routers = append(srv.routers, r)
}

func (srv *Server) AddStaticRouter(match, prefix, filePath string) {
	re, err := regexp.Compile(match)
	if err != nil {
		srv.log.Error("compile regexp error: %v", err)
		return
	}
	fs := http.FileServer(http.Dir(filePath))
	handler := http.StripPrefix(prefix, fs)
	r := &staticRouter{
		baseRouter: &baseRouter{
			re: re,
		},
		handler: handler,
	}
	srv.routers = append(srv.routers, r)

}

func (srv *Server) setSession(ctx *Context) {
	ctx.session = srv.sm.session(ctx)
}
