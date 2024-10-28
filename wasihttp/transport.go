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

// Transport implements [http.RoundTripper].
type Transport struct{}

// RoundTrip executes a single HTTP transaction, using [wasi-http] APIs.
//
// [wasi-http]: https://github.com/webassembly/wasi-http
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	outgoingRequest := types.NewOutgoingRequest(toFields(req.Header))
	outgoingRequest.SetAuthority(cm.Some(requestAuthority(req))) // TODO: when should this be cm.None?
	outgoingRequest.SetMethod(toMethod(req.Method))
	outgoingRequest.SetPathWithQuery(requestPath(req))
	outgoingRequest.SetScheme(cm.Some(toScheme(req.URL.Scheme))) // TODO: when should this be cm.None?

	outgoingBody_ := outgoingRequest.Body()
	outgoingBody := *outgoingBody_.OK() // the first call should always return OK

	// TODO: when are [options] used?
	// [options]: https://github.com/WebAssembly/wasi-http/blob/main/wit/handler.wit#L38-L39
	futureIncomingResponse_ := outgoinghandler.Handle(outgoingRequest, cm.None[types.RequestOptions]())
	if err := checkError(futureIncomingResponse_); err != nil {
		// outgoing request is invalid or not allowed to be made
		return nil, err
	}

	writeOutgoingBody(req.Body, outgoingBody)

	// Finalize the request body
	// TODO: complete the request trailers
	finish := types.OutgoingBodyFinish(outgoingBody, cm.None[types.Fields]())
	if err := checkError(finish); err != nil {
		return nil, err
	}

	futureIncomingResponse := futureIncomingResponse_.OK()
	defer futureIncomingResponse.ResourceDrop()
	poll := futureIncomingResponse.Subscribe()
	defer poll.ResourceDrop()
	if !poll.Ready() {
		poll.Block()
	}
	someResponse := futureIncomingResponse.Get()
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

// writeOutgoingBody writes the io.ReadCloser to the wasi-http [types.OutgoingBody].
//
// [types.writeOutgoingBody]: https://github.com/WebAssembly/wasi-http/blob/v0.2.0/wit/types.wit#L514-L540
func writeOutgoingBody(body io.ReadCloser, wasiBody types.OutgoingBody) error {
	defer body.Close()
	w := bodyWriter(wasiBody)
	defer w.Close()
	if _, err := io.Copy(w, body); err != nil {
		return fmt.Errorf("wasihttp: %v", err)
	}
	return nil
}

func checkError[Shape, OK, Err any](result cm.Result[Shape, OK, Err]) error {
	if result.IsErr() {
		return fmt.Errorf("wasihttp: %v", result.Err())
	}
	return nil
}
