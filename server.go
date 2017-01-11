package goweb

import (
	"net/http"
	"regexp"

	"github.com/dafengge0913/golog"
	"net"
)

const defaultSessionTimeout = 600 // 10 minute

type HTTPHandler func(ctx *Context)

type Server struct {
	ln            net.Listener
	log           *golog.Logger
	routers       []*router
	staticRouters map[string]http.Handler
	sm            *sessionManager
}

func NewServer() *Server {
	return &Server{
		log:           golog.NewLogger(golog.LEVEL_DEBUG, nil),
		routers:       make([]*router, 0),
		staticRouters: make(map[string]http.Handler),
		sm:            newSessionManager(defaultSessionTimeout),
	}
}

type router struct {
	re      *regexp.Regexp
	handler HTTPHandler
}

// start http Server
// call after all routers have been added
func (srv *Server) Serve(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv.ln = ln
	defer ln.Close()
	handler := http.NewServeMux()
	handler.Handle("/", srv)
	for m, h := range srv.staticRouters {
		handler.Handle(m, h)
	}
	srv.staticRouters = nil // staticRouters won't be used anymore
	http.Serve(ln, handler)
	return nil
}

func (srv *Server) Close() {
	if srv.ln != nil {
		if err := srv.ln.Close(); err != nil {
			srv.log.Info("close server error: %v", err)
		}
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
	r.handler(ctx)
}

func (srv *Server) findRouter(path string) *router {
	var mostMatch *router
	mostMatchIdx := 0
	for _, r := range srv.routers {
		if !r.re.MatchString(path) {
			continue
		}
		match := r.re.FindStringSubmatchIndex(path)
		matchLen := len(match)
		if matchLen == 0 {
			continue
		}
		matchIdx := match[matchLen-1]
		if matchIdx >= mostMatchIdx { // use the last setting
			mostMatchIdx = matchIdx
			mostMatch = r
		}

	}
	return mostMatch
}

func (srv *Server) AddRouter(match string, handler HTTPHandler) {
	matchRe, err := regexp.Compile(match)
	if err != nil {
		srv.log.Error("compile regexp error: %v", err)
		return
	}
	r := &router{
		re:      matchRe,
		handler: handler,
	}
	srv.routers = append(srv.routers, r)
}

func (srv *Server) AddStaticRouter(match, filePath string) {
	fs := http.FileServer(http.Dir(filePath))
	srv.staticRouters[match] = http.StripPrefix(match, fs)
}

func (srv *Server) setSession(ctx *Context) {
	ctx.session = srv.sm.getOrCreateSession(ctx)
}
