package main

import (
	"net/http"

	_ "github.com/ydnar/wasi-http-go/wasihttp"
)

func init() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.String())
	})

	http.HandleFunc("GET /safe", func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.String())
	})
}

func main() {}
