# wasi-http-go

[![pkg.go.dev](https://img.shields.io/badge/docs-pkg.go.dev-blue.svg)](https://pkg.go.dev/github.com/ydnar/wasi-http-go) [![build status](https://img.shields.io/github/actions/workflow/status/ydnar/wasi-http-go/test.yaml?branch=main)](https://github.com/ydnar/wasi-http-go/actions)

## [wasi-http](https://github.com/WebAssembly/wasi-http) for [Go](https://go.dev)

Package `wasihttp` implements the [`wasi:http/proxy`](https://github.com/WebAssembly/wasi-http/blob/v0.2.0/proxy.md) version 0.2.0 for Go using standard [`net/http`](https://pkg.go.dev/net/http) interfaces.

## Prerequisites

To use this package, you’ll need:

- [TinyGo](https://tinygo.org/) 0.34.0 or later.
- [Wasmtime](https://wasmtime.dev/) 26.0.0 or later.

## Examples

Example code using this package can be found in the [examples](./examples) directory. To run the examples with `tinygo run`, you’ll need to install a development build of TinyGo that supports `wasmtime serve` (0.35.0-dev or later with [this PR](https://github.com/tinygo-org/tinygo/pull/4555) merged).

### Server

A simple `net/http` server. Run with `tinygo run -target wasip2-http.json ./main.go`.

```go
package main

import (
	"net/http"
	_ "github.com/ydnar/wasi-http-go/wasihttp" // enable wasi-http
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Go", "Gopher")
		w.Write([]byte("Hello world!\n"))
	})
}

func main() {}
```

## License

This project is licensed under the Apache 2.0 license with the LLVM exception. See [LICENSE](LICENSE) for more details.
