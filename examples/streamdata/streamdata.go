// This example is taken from https://github.com/dev-wasm/dev-wasm-go/blob/main/http/main.go
// demonstrates how to use the wasihttp package to make HTTP requests using the `http.Client` interface.
//
// To run: `tinygo build -target=wasip2-roundtrip.json -o streamdata.wasm ./examples/streamdata`
// Test: `wasmtime run -Shttp -Sinherit-network -Sinherit-env streamdata.wasm`
package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"

	wasihttp "github.com/ydnar/wasi-http-go/wasihttp"
)

//go:embed bigdata_10kb.data
var bigdata []byte

func printResponse(r *http.Response, printBody bool) error {
	fmt.Printf("Status: %d\n", r.StatusCode)
	for k, v := range r.Header {
		fmt.Printf("%s: %s\n", k, v[0])
	}
	if r.Body != nil && printBody {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		fmt.Printf("Body: \n%s\n", body)
	}

	return nil
}

func main() {
	client := &http.Client{
		Transport: &wasihttp.Transport{},
	}
	req, err := http.NewRequest(http.MethodPost, "https://postman-echo.com/post", bytes.NewReader(bigdata))
	if err != nil {
		panic(err.Error())
	}
	if req == nil {
		panic("Nil request!")
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(bigdata)))
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	err = printResponse(res, false)
	if err != nil {
		panic(err.Error())
	}
}
