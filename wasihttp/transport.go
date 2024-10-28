package wasihttp

import (
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
	handled := outgoinghandler.Handle(r, cm.None[types.RequestOptions]())
	if err := checkError(handled); err != nil {
		// outgoing request is invalid or not allowed to be made
		return nil, err
	}
	incoming := handled.OK()
	defer incoming.ResourceDrop()

	somebody := r.Body()
	body := *somebody.OK() // the first call should always return OK

	// Write request body
	w := bodyWriter(body)
	defer w.finish()
	if _, err := io.Copy(w, req.Body); err != nil {
		return nil, fmt.Errorf("wasihttp: %v", err)
	}
	w.Flush()

	// Finalize the request body
	// TODO: complete the request trailers
	finished := types.OutgoingBodyFinish(body, cm.None[types.Fields]())
	if err := checkError(finished); err != nil {
		return nil, err
	}

	// Wait for response
	poll := incoming.Subscribe()
	if !poll.Ready() {
		poll.Block()
	}
	poll.ResourceDrop()

	someResponse := incoming.Get()
	if someResponse.None() {
		return nil, fmt.Errorf("wasihttp: future response is None after blocking")
	}

	responseResult := someResponse.Some().OK() // the first call should always return OK
	if err := checkError(*responseResult); err != nil {
		// TODO: what do we do with the HTTP proxy error-code?
		return nil, err
	}

	response := *responseResult.OK()
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

func checkError[Shape, OK, Err any](result cm.Result[Shape, OK, Err]) error {
	if result.IsErr() {
		return fmt.Errorf("wasihttp: %v", result.Err())
	}
	return nil
}
