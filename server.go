package goweb

import (
	"github.com/dafengge0913/golog"
	"net/http"
	"regexp"
)

type HTTPHandler func(ctx *Context)

type server struct {
	log     *golog.Logger
	routers []*router
}

func NewServer() *server {
	return &server{
		log:     golog.NewLogger(golog.LEVEL_DEBUG, nil),
		routers: make([]*router, 0),
	}

}

func (srv *server) Serve(addr string) error {
	handler := http.NewServeMux()
	handler.Handle("/", srv)
	return http.ListenAndServe(addr, handler)
}

func (srv *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	if err := req.ParseForm(); err != nil {
		srv.log.Error("Parse request form error: %v", err)
		return
	}
	ctx := newContext(srv.log, req, rw)
	req.ParseForm()
	if len(req.Form) > 0 {
		for k, v := range req.Form {
			ctx.Params[k] = v[0]
		}
	}
	r := srv.findRouter(reqPath)
	if r == nil {
		srv.log.Debug("router not found path:%s", reqPath)
		return
	}
	r.handler(ctx)
}

func (srv *server) findRouter(path string) *router {
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
		matchIdx := match[matchLen - 1]
		if matchIdx >= mostMatchIdx { // use the last setting
			mostMatchIdx = matchIdx
			mostMatch = r
		}

	}
	return mostMatch
}

func (srv *server) AddRouter(match string, handler HTTPHandler) {
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
