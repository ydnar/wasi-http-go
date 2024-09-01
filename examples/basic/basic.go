package main

// import (
// 	incominghandler "github.com/ydnar/wasi-http-go/internal/wasi/http/incoming-handler"
// 	"github.com/ydnar/wasi-http-go/internal/wasi/http/types"
// )

import (
	_ "os"
	_ "syscall"
	_ "unsafe"

	_ "github.com/ydnar/wasi-http-go/wasihttp"
	"github.com/ydnar/wasm-tools-go/cm"
)

func main() {}

// func init() {
// 	// Assign the "wasi:http/incoming-handler@0.2.1#handle" export.
// 	incominghandler.Exports.Handle = handleIncomingRequest
// }

// func handleIncomingRequest(req types.IncomingRequest, res types.ResponseOutparam) {
// 	panic("wasi:http/incoming-handler@0.2.1#handle")
// }

//export wasi:cli/environment@0.2.0#get-environment
func wasmimport_GetEnvironment(result *cm.List[[2]string]) {}

//export wasi:cli/environment@0.2.0#get-arguments
func wasmimport_GetArguments(result *cm.List[string]) {}

//export wasi:cli/environment@0.2.0#initial-cwd
func wasmimport_InitialCWD(result *cm.Option[string]) {}
