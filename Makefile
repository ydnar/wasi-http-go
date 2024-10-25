.PHONY: tools
tools:
	go generate -tags tools ./...

.PHONY: go-bindings
go-bindings:
	go run github.com/bytecodealliance/wasm-tools-go/cmd/wit-bindgen-go generate -o internal/ ./wit
