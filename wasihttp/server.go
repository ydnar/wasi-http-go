package wasihttp

import (
	"io"
	"net/http"
	"net/url"
	"strings"

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
	w := newResponseWriter(req, res)
	h.ServeHTTP(w, w.req)
}

var _ http.ResponseWriter = &responseWriter{}

type responseWriter struct {
	outparam    types.ResponseOutparam
	req         *http.Request
	header      http.Header
	wroteHeader bool
	status      int // HTTP status code passed to WriteHeader

	res types.OutgoingResponse // valid after outparam is set
}

func newResponseWriter(req types.IncomingRequest, res types.ResponseOutparam) *responseWriter {
	return &responseWriter{
		outparam: res,
		req:      incomingRequest(req),
		header:   make(http.Header),
	}
}

func (w *responseWriter) Header() http.Header {
	// TODO: handle concurrent access to (or mutations of) w.header?
	return w.header
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return 0, nil
}

func (w *responseWriter) WriteHeader(code int) {
	println("responseWriter.WriteHeader")
	if w.wroteHeader {
		// TODO: improve logging
		println("already wrote header")
		return
	}

	// TODO: handle 1xx informational headers?

	w.wroteHeader = true
	w.status = code

	headers := toFields(w.header)
	res := types.NewOutgoingResponse(headers)
	res.SetStatusCode(types.StatusCode(code))

	types.ResponseOutparamSet(w.outparam, cm.OK[outgoingResult](res))

	// TODO
}

type outgoingResult = cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]

// fatal sets an error code on the response, to allow the implementation
// to determine how to respond with an HTTP error response.
func (w *responseWriter) fatal(e types.ErrorCode) {
	types.ResponseOutparamSet(w.outparam, cm.Err[outgoingResult](e))
}

func incomingRequest(req types.IncomingRequest) *http.Request {
	var body io.ReadCloser
	rbody := req.Consume()
	if b := rbody.OK(); b != nil {
		body = &incomingBody{*b}
	}
	r := &http.Request{
		Method: method(req.Method()),
		URL:    incomingURL(req),
		// TODO: Proto, ProtoMajor, ProtoMinor
		Header: httpHeader(req.Headers()),
		Host:   req.Authority().Value(),
		Body:   body,
	}
	return r
}

var _ io.ReadCloser = &incomingBody{}

type incomingBody struct {
	types.IncomingBody
}

func (b *incomingBody) Read(p []byte) (int, error) {
	// TODO
	return 0, nil
}

func (b *incomingBody) Close() error {
	// TODO
	return nil
}

func method(m types.Method) string {
	if o := m.Other(); o != nil {
		return strings.ToUpper(*o)
	}
	return strings.ToUpper(m.String())
}

func incomingURL(req types.IncomingRequest) *url.URL {
	u := &url.URL{
		Scheme: scheme(req.Scheme().Value()),
		Host:   req.Authority().Value(),
	}
	u, _ = u.Parse(req.PathWithQuery().Value())
	return u
}

func scheme(s types.Scheme) string {
	if o := s.Other(); o != nil {
		return strings.ToLower(*o)
	}
	return strings.ToLower(s.String())
}

func httpHeader(fields types.Fields) http.Header {
	h := http.Header{}
	for _, e := range fields.Entries().Slice() {
		h.Add(string(e.F0), string(e.F1.Slice()))
	}
	return h
}
