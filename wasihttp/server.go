package wasihttp

import (
	"net/http"

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
	panic("wasi:http/incoming-handler@0.2.1#handle")
}
