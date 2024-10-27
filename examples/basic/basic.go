package main

import (
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
}

func main() {}
