package wasihttp

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	outgoinghandler "github.com/ydnar/wasi-http-go/internal/wasi/http/outgoing-handler"
	"github.com/ydnar/wasi-http-go/internal/wasi/http/types"
)

var _ http.RoundTripper = &Transport{}

type Transport struct{}

// RoundTrip executes a single HTTP transaction, using [wasi-http] APIs.
//
// [wasi-http]: https://github.com/webassembly/wasi-http
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	defer req.Body.Close()

	// TODO: wrap this into a helper func outgoingRequest?
	r := types.NewOutgoingRequest(toFields(req.Header))
	r.SetAuthority(cm.Some(requestAuthority(req))) // TODO: when should this be cm.None?
	r.SetMethod(toMethod(req.Method))
	r.SetPathWithQuery(requestPath(req))
	r.SetScheme(cm.Some(toScheme(req.URL.Scheme))) // TODO: when should this be cm.None?

	// TODO: when are [options] used?
	// [options]: https://github.com/WebAssembly/wasi-http/blob/main/wit/handler.wit#L38-L39
	incoming, err, isErr := outgoinghandler.Handle(r, cm.None[types.RequestOptions]()).Result()
	if isErr {
		// outgoing request is invalid or not allowed to be made
		return nil, errors.New(err.String())
	}
	defer incoming.ResourceDrop()

	body, _, _ := r.Body().Result() // the first call should always return OK

	// Write request body
	w := newBodyWriter(body, func() http.Header {
		// TODO: extract request trailers
		return nil
	})
	if _, err := io.Copy(w, req.Body); err != nil {
		return nil, fmt.Errorf("wasihttp: %v", err)
	}
	w.Flush()
	w.finish()

	// Wait for response
	poll := incoming.Subscribe()
	if !poll.Ready() {
		poll.Block()
	}
	poll.ResourceDrop()

	future := incoming.Get()
	if future.None() {
		return nil, fmt.Errorf("wasihttp: future response is None after blocking")
	}
	// TODO: figure out a better way to handle option<result<result<incoming-response, error-code>>>
	response, err, isErr := future.Some().OK().Result() // the first call should always return OK
	if isErr {
		// TODO: what do we do with the HTTP proxy error-code?
		return nil, errors.New(err.String())
	}
	defer response.ResourceDrop()

	return incomingResponse(response)
}

func requestAuthority(req *http.Request) string {
	if req.Host == "" {
		return req.URL.Host
	} else {
		return req.Host
	}
}

func requestPath(req *http.Request) cm.Option[string] {
	path := req.URL.RequestURI()
	if path == "" {
		return cm.None[string]()
	}
	return cm.Some(path)
}
