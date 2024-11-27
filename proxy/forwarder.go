package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

type errorHandlerFn func(http.ResponseWriter, *http.Request, error)

type Forwarder struct {
	cfg     ForwarderConfig
	proxy   *httputil.ReverseProxy
	target  *url.URL
	tracker *UpstreamTracker
}

func NewForwarder(cfg ForwarderConfig) *Forwarder {
	target, err := url.Parse(cfg.Upstream)
	if err != nil {
		zap.L().Fatal("unable to parse upstream gateway URL", zap.Error(err))
	}

	tracker := NewUpstreamTracker()
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		MaxIdleConns:          cfg.MaxIdleConns,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		DisableCompression:    cfg.DisableCompression,
		ResponseHeaderTimeout: cfg.HeaderTimeout,
	}

	return &Forwarder{cfg, proxy, target, tracker}
}

func (f *Forwarder) Forward(piece string, w http.ResponseWriter, r *http.Request) {
	originalDirector := f.proxy.Director
	f.proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.RawQuery = ""
		req.URL.Fragment = ""
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Origin-Host", f.target.Host)
		req.URL.Path = fmt.Sprintf("/piece/%s", piece)
	}

	f.proxy.ServeHTTP(w, r)
}

func (f *Forwarder) SetErrorHandler(handler errorHandlerFn) {
	f.proxy.ErrorHandler = handler
}

func (f *Forwarder) SetResponseTransformer(handler func(*http.Response) error) {
	f.proxy.ModifyResponse = handler
}
