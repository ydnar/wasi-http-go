package wasihttp

import (
	"net/http"

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

func handleIncomingRequest(req types.IncomingRequest, out types.ResponseOutparam) {
	h := defaultHandler
	if h == nil {
		h = http.DefaultServeMux
	}
	w, err := newResponseWriter(req, out)
	if err != nil {
		return // TODO: log error?
	}
	h.ServeHTTP(w, w.req)
}

var _ http.ResponseWriter = &responseWriter{}

type responseWriter struct {
	out         types.ResponseOutparam
	req         *http.Request
	header      http.Header
	wroteHeader bool
	status      int // HTTP status code passed to WriteHeader

	res types.OutgoingResponse // valid after outparam is set
}

func newResponseWriter(req types.IncomingRequest, resout types.ResponseOutparam) (*responseWriter, error) {
	r, err := incomingRequest(req)
	w := &responseWriter{
		out:    resout,
		req:    r,
		header: make(http.Header),
	}
	if err != nil {
		w.fatal(types.ErrorCodeHTTPProtocolError())
	}
	return w, err
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

	types.ResponseOutparamSet(w.out, cm.OK[outgoingResult](res))

	// TODO
}

type outgoingResult = cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]

// fatal sets an error code on the response, to allow the implementation
// to determine how to respond with an HTTP error response.
func (w *responseWriter) fatal(e types.ErrorCode) {
	types.ResponseOutparamSet(w.out, cm.Err[outgoingResult](e))
}
