// This example implements a basic web server.
//
// To run: `tinygo run -target=wasip2-http.json ./examples/basic`
// Test /: `curl -v 'http://0.0.0.0:8080/'`
// Test /safe: `curl -v 'http://0.0.0.0:8080/safe'`
// Test /counter: `curl -v 'http://0.0.0.0:8080/counter'`

package main

import (
	"fmt"
	"net/http"

	_ "github.com/ydnar/wasi-http-go/wasihttp"
)

func init() {
	// TODO: use "GET /" when TinyGo supports net/http from Go 1.22+
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Go", "Gopher")
		w.Write([]byte("Hello world!\n"))
	})

	http.HandleFunc("/safe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Go", "Safe")
		w.Write([]byte("Welcome to /safe\n"))
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		// do nothing, force default response handling
	})

	http.HandleFunc("/counter", func(w http.ResponseWriter, r *http.Request) {
		n := <-c
		fmt.Fprintf(w, "%d", n)
	})

	go func() {
		var n int
		for {
			c <- n
			n++
		}
	}()
}

var c = make(chan int, 1)

func main() {}
