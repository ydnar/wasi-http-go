.PHONY: go-bindings
go-bindings:
	wit-bindgen-go generate -o internal/ --world go:http/proxy ./wit
