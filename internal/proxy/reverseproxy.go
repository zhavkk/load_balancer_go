package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxyAdapter struct {
	proxy *httputil.ReverseProxy
}

type Proxy interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	SetErrorHandler(fn func(http.ResponseWriter, *http.Request, error))
}

func NewReverseProxy(target *url.URL) Proxy {
	p := httputil.NewSingleHostReverseProxy(target)
	return &ReverseProxyAdapter{proxy: p}
}

func (r *ReverseProxyAdapter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.proxy.ServeHTTP(w, req)
}

func (r *ReverseProxyAdapter) SetErrorHandler(fn func(http.ResponseWriter, *http.Request, error)) {
	r.proxy.ErrorHandler = fn
}
