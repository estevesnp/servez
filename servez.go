package servez

import (
	"fmt"
	"log"
	"net/http"

	m "github.com/estevesnp/servez/middleware"
)

type Router struct {
	handler             *http.ServeMux
	addr                string
	preMiddlewareFuncs  []http.HandlerFunc
	postMiddlewareFuncs []http.HandlerFunc
}

type RouterCfg struct {
	Addr                string
	PreMiddlewareFuncs  []http.HandlerFunc
	PostMiddlewareFuncs []http.HandlerFunc
}

var defaultCfg = RouterCfg{
	Addr:                "localhost:8080",
	PreMiddlewareFuncs:  []http.HandlerFunc{m.LogRequest},
	PostMiddlewareFuncs: []http.HandlerFunc{m.LogResponse},
}

type httpVerb string

const (
	getVerb    = "GET"
	postVerb   = "POST"
	putVerb    = "PUT"
	deleteVerb = "DELETE"
)

func New(routerCfg *RouterCfg) *Router {
	cfg := defaultCfg

	if routerCfg != nil {
		if routerCfg.Addr != "" {
			cfg.Addr = routerCfg.Addr
		}

		if routerCfg.PreMiddlewareFuncs != nil {
			cfg.PreMiddlewareFuncs = routerCfg.PreMiddlewareFuncs
		}

		if routerCfg.PostMiddlewareFuncs != nil {
			cfg.PostMiddlewareFuncs = routerCfg.PostMiddlewareFuncs
		}
	}

	return &Router{
		handler:             http.NewServeMux(),
		addr:                cfg.Addr,
		preMiddlewareFuncs:  cfg.PreMiddlewareFuncs,
		postMiddlewareFuncs: cfg.PostMiddlewareFuncs,
	}
}

func Default() *Router {
	return New(nil)
}

func applyPreMiddleware(middlewareFunc, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middlewareFunc(w, r)
		handler(w, r)
	}
}

func applyPostMiddleware(middlewareFunc, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
		middlewareFunc(w, r)
	}
}

func (r *Router) genericHandler(verb httpVerb, pattern string, handler http.HandlerFunc) {
	endpoint := fmt.Sprintf("%s %s", verb, pattern)

	for i := len(r.preMiddlewareFuncs) - 1; i >= 0; i-- {
		f := r.preMiddlewareFuncs[i]
		handler = applyPreMiddleware(f, handler)
	}

	for _, f := range r.postMiddlewareFuncs {
		handler = applyPostMiddleware(f, handler)
	}

	r.handler.HandleFunc(endpoint, handler)
}

func (r *Router) GET(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.genericHandler(getVerb, pattern, handler)
}

func (r *Router) POST(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.genericHandler(postVerb, pattern, handler)
}

func (r *Router) PUT(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.genericHandler(putVerb, pattern, handler)
}

func (r *Router) DELETE(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.genericHandler(deleteVerb, pattern, handler)
}

func (r *Router) Start() error {
	log.Printf("Starting server on %s\n", r.addr)
	return http.ListenAndServe(r.addr, r.handler)
}
