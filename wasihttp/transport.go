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

	outgoingBody_ := outgoingRequest.Body()
	outgoingBody := outgoingBody_.OK() // the first call should always return OK

	// TODO: when are [options] used?
	// [options]: https://github.com/WebAssembly/wasi-http/blob/main/wit/handler.wit#L38-L39
	futureIncomingResponse_ := outgoinghandler.Handle(outgoingRequest, cm.None[types.RequestOptions]())

	if err := checkError(futureIncomingResponse_); err != nil {
		// outgoing request is invalid or not allowed to be made
		return nil, err
	}

	toBody(&req.Body, outgoingBody)

	// Finalize the request body
	// TODO: complete the request trailers
	finish := types.OutgoingBodyFinish(*outgoingBody, cm.None[types.Fields]())
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
	responseBody := futureIncomingResponse.Get()
	if responseBody.None() {
		return nil, fmt.Errorf("wasihttp: response is None after blocking")
	}

	incomingResponse_ := responseBody.Some().OK() // the first call should always return OK
	if err := checkError(*incomingResponse_); err != nil {
		// TODO: what do we do with the HTTP proxy error-code?
		return nil, err
	}
	incomingResponse := incomingResponse_.OK()
	defer incomingResponse.ResourceDrop()

	response := &http.Response{
		Status:     http.StatusText(int(incomingResponse.Status())),
		StatusCode: int(incomingResponse.Status()),
		Header:     FromHeaders(incomingResponse.Headers()),
	}

	ib := incomingResponse.Consume()
	if err := checkError(ib); err != nil {
		return nil, err
	}
	incomingBody := ib.OK()

	response.Body = fromBody(incomingBody)
	return response, nil
}

func checkError[Shape, Ok, Err any](result cm.Result[Shape, Ok, Err]) error {
	if result.IsErr() {
		return fmt.Errorf("wasihttp: %v", result.Err())
	}
	return nil
}
