package wasihttp

import (
	"net/http"
	"net/url"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	incominghandler "github.com/ydnar/wasi-http-go/internal/wasi/http/incoming-handler"
	"github.com/ydnar/wasi-http-go/internal/wasi/http/types"
)

var defaultHandler http.Handler

// Serve sets the [http.Handler] that incoming [wasi-http] requests are routed to.
// If not set, [http.DefaultServeMux] is used.
func Serve(h http.Handler) {
	defaultHandler = h
}

func init() {
	// Assign the "wasi:http/incoming-handler@0.2.1#handle" export.
	incominghandler.Exports.Handle = handleIncomingRequest
}

func handleIncomingRequest(req types.IncomingRequest, res types.ResponseOutparam) {
	h := defaultHandler
	if h == nil {
		h = http.DefaultServeMux
	}
	w := http.ResponseWriter(nil)
	r := incomingRequest(req)
	h.ServeHTTP(w, r)
}

func incomingRequest(req types.IncomingRequest) *http.Request {
	r := &http.Request{
		Method: method(req.Method()),
		URL:    incomingURL(req),
		// TODO: Proto, ProtoMajor, ProtoMinor
		Header: header(req.Headers()),
		Host:   optionZero(req.Authority()),
	}
	return r
}

func optionZero[T any](o cm.Option[T]) T {
	if o.None() {
		var zero T
		return zero
	}
	return *o.Some()
}

func method(m types.Method) string {
	switch {
	case m.Connect():
		return "CONNECT"
	case m.Delete():
		return "DELETE"
	case m.Get():
		return "GET"
	case m.Head():
		return "HEAD"
	case m.Options():
		return "OPTIONS"
	case m.Patch():
		return "PATCH"
	case m.Post():
		return "POST"
	case m.Put():
		return "PUT"
	case m.Trace():
		return "TRACE"
	}
	if o := m.Other(); o != nil {
		return *o
	}
	return ""
}

func incomingURL(req types.IncomingRequest) *url.URL {
	return &url.URL{
		Scheme: scheme(req.Scheme()),
	}
}

func scheme(o cm.Option[types.Scheme]) string {
	if o.None() {
		return ""
	}
	s := *o.Some()
	switch {
	case s.HTTP():
		return "http"
	case s.HTTPS():
		return "https"
	}
	if other := s.Other(); other != nil {
		return *other
	}
	return ""
}

func header(fields types.Fields) http.Header {
	h := http.Header{}
	for _, e := range fields.Entries().Slice() {
		h.Add(string(e.F0), string(e.F1.Slice()))
	}
	return h
}
