package wasihttp

import "net/http"

var _ http.RoundTripper = &Transport{}

// Transport implements [http.RoundTripper].
type Transport struct{}

// RoundTrip executes a single HTTP transaction, using [wasi-http] APIs.
//
// [wasi-http]: https://github.com/webassembly/wasi-http
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, nil
}
