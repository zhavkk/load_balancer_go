package proxy

import "net/http"

type ProxyHandler struct {
	Handler http.Handler
}

func NewProxyHandler(h http.Handler) *ProxyHandler {
	return &ProxyHandler{Handler: h}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handler.ServeHTTP(w, r)
}
