package wasihttp

import (
	"errors"
	"fmt"
	"net/http"
	"os"

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
	w.finish()
}

var _ http.ResponseWriter = &responseWriter{}

type responseWriter struct {
	out         types.ResponseOutparam
	req         *http.Request
	header      http.Header
	wroteHeader bool
	status      int // HTTP status code passed to WriteHeader

	res    types.OutgoingResponse // valid after headers are sent
	body   types.OutgoingBody     // valid after res.Body() is called
	writer *bodyWriter            // valid after body.Stream() is called

	finished bool
}

func newResponseWriter(req types.IncomingRequest, out types.ResponseOutparam) (*responseWriter, error) {
	r, err := incomingRequest(req)
	w := &responseWriter{
		out:    out,
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
	if w.finished {
		return 0, errors.New("wasihttp: write after close")
	}
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.writer.Write(p)
}

func (w *responseWriter) WriteHeader(code int) {
	if w.finished || w.wroteHeader {
		// TODO: improve logging
		return
	}

	// TODO: handle 1xx informational headers?

	w.wroteHeader = true
	w.status = code

	headers := toFields(w.header)
	w.res = types.NewOutgoingResponse(headers)
	w.res.SetStatusCode(types.StatusCode(code))

	rbody := w.res.Body()
	w.body = *rbody.OK() // the first call should always return OK
	w.writer = newBodyWriter(w.body)

	// Consume the response-outparam and outgoing-response.
	types.ResponseOutparamSet(w.out, cm.OK[outgoingResult](w.res))
}

func (w *responseWriter) finish() {
	if w.finished {
		return
	}
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	w.finished = true
	w.writer.finish()

	// TODO: extract trailers from http.ResponseWriter
	var trailers cm.Option[types.Trailers]

	result := types.OutgoingBodyFinish(w.body, trailers)
	if result.IsErr() {
		// TODO: improve this
		fmt.Fprintf(os.Stderr, "wasihttp: outgoing-body-finish: %v", result.Err())
	}
}

type outgoingResult = cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]

// fatal sets an error code on the response, to allow the implementation
// to determine how to respond with an HTTP error response.
func (w *responseWriter) fatal(e types.ErrorCode) {
	w.finished = true
	types.ResponseOutparamSet(w.out, cm.Err[outgoingResult](e))
}
