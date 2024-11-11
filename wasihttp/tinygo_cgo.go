//go:build tinygo.wasm

package wasihttp

// #cgo LDFLAGS: -L. -lgowasihttp-embed
import "C"
