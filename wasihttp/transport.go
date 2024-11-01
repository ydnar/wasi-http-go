package wasihttp

import (
	"fmt"
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
	hdrs := toFields(req.Header)

	outgoingRequest := types.NewOutgoingRequest(hdrs)

	auth := authority(req)
	// TODO: when should we set the authority to `cm.None`?
	outgoingRequest.SetAuthority(cm.Some(auth))

	m := toMethod(req.Method)
	outgoingRequest.SetMethod(m)

	p := path(req)
	outgoingRequest.SetPathWithQuery(p)

	scheme := toScheme(req.URL.Scheme)
	outgoingRequest.SetScheme(cm.Some(scheme))

	outgoingBody, _, _ := outgoingRequest.Body().Result() // the first call should always return OK

	// TODO: when are [options] used?
	// [options]: https://github.com/WebAssembly/wasi-http/blob/main/wit/handler.wit#L38-L39
	futureIncomingResponse, err, isErr := outgoinghandler.Handle(outgoingRequest, cm.None[types.RequestOptions]()).Result()
	defer futureIncomingResponse.ResourceDrop()

	if isErr {
		// outgoing request is invalid or not allowed to be made
		return nil, fmt.Errorf("wasihttp: %v", err)
	}

	toBody(&req.Body, outgoingBody)

	// Finalize the request body
	// TODO: complete the request trailers
	_, err, isErr = types.OutgoingBodyFinish(outgoingBody, cm.None[types.Fields]()).Result()
	if isErr {
		return nil, fmt.Errorf("wasihttp: %v", err)
	}

	poll := futureIncomingResponse.Subscribe()
	defer poll.ResourceDrop()
	if !poll.Ready() {
		poll.Block()
	}
	responseBody := futureIncomingResponse.Get()
	if responseBody.None() {
		return nil, fmt.Errorf("wasihttp: response is None after blocking")
	}

	incomingResponse, err, isErr := responseBody.Some().OK().Result() // the first call should always return OK
	defer incomingResponse.ResourceDrop()
	if isErr {
		// TODO: what do we do with the HTTP proxy error-code?
		return nil, fmt.Errorf("wasihttp: %v", err)
	}

	response := &http.Response{
		Status:     http.StatusText(int(incomingResponse.Status())),
		StatusCode: int(incomingResponse.Status()),
		Header:     FromHeaders(incomingResponse.Headers()),
	}

	incomingBody, err_, isErr := incomingResponse.Consume().Result()
	if isErr {
		return nil, fmt.Errorf("wasihttp: %v", err_)
	}

	response.Body = fromBody(incomingBody)
	return response, nil
}
