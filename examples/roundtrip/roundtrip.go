package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	wasihttp "github.com/ydnar/wasi-http-go/wasihttp"
)

func printResponse(r *http.Response) {
	fmt.Printf("Status: %d\n", r.StatusCode)
	for k, v := range r.Header {
		fmt.Printf("%s: %s\n", k, v[0])
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Body: \n%s\n", body)
}

func main() {
	client := &http.Client{
		Transport: &wasihttp.Transport{},
	}
	req, err := http.NewRequest("GET", "https://postman-echo.com/get", nil)
	if err != nil {
		panic(err.Error())
	}
	if req == nil {
		panic("Nil request!")
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()
	printResponse(res)

	res, err = client.Post("https://postman-echo.com/post", "application/json", bytes.NewReader([]byte("{\"foo\": \"bar\"}")))
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()
	printResponse(res)

	req, err = http.NewRequest("PUT", "http://postman-echo.com/put", bytes.NewReader([]byte("{\"baz\": \"blah\"}")))
	if err != nil {
		panic(err.Error())
	}
	if req == nil {
		panic("Nil request!")
	}
	res, err = client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()
	printResponse(res)
}
