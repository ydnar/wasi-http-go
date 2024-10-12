package wasihttp

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/ydnar/wasi-http-go/internal/wasi/http/types"
	"github.com/ydnar/wasi-http-go/internal/wasi/io/streams"
)

// authority returns the authority of the request from the [http.Request] Host or URL.Host.
//
// Note the description of the Host field: https://cs.opensource.google/go/go/+/refs/tags/go1.23.2:src/net/http/request.go;l=240-243
func authority(req *http.Request) string {
	if req.Host == "" {
		return req.URL.Host
	} else {
		return req.Host
	}
}

// toFields convert the [http.Header] to a wasi-http [types.Fields].
//
// [types.Fields]: https://github.com/WebAssembly/wasi-http/blob/v0.2.0/wit/types.wit#L215
func toFields(headers http.Header) types.Fields {
	fields := types.NewFields()
	for k, v := range headers {
		key := types.FieldKey(k)
		vals := []types.FieldValue{}
		for _, vv := range v {
			vals = append(vals, types.FieldValue(cm.ToList([]uint8(vv))))
		}
		fields.Set(key, cm.ToList(vals))
	}
	return fields
}

// FromHeaders convert the wasi-http [types.Fields] to a [http.Header].
//
// [types.Fields]: https://github.com/WebAssembly/wasi-http/blob/v0.2.0/wit/types.wit#L215
func FromHeaders(fields types.Fields) http.Header {
	h := http.Header{}
	es := fields.Entries()
	for _, field := range es.Slice() {
		k, v := string(field.F0), string(field.F1.Slice())
		h.Add(k, v)
	}
	fields.ResourceDrop()
	return h
}

// toMethod returns the wasi-http [types.toMethod] from the [http.Request] method.
// If the method is not a standard HTTP method, it returns a [types.MethodOther].
//
// [types.toMethod]: https://github.com/WebAssembly/wasi-http/blob/v0.2.0/wit/types.wit#L11-L22
func toMethod(s string) types.Method {
	switch s {
	case http.MethodGet:
		return types.MethodGet()
	case http.MethodHead:
		return types.MethodHead()
	case http.MethodPost:
		return types.MethodPost()
	case http.MethodPut:
		return types.MethodPut()
	case http.MethodPatch:
		return types.MethodPatch()
	case http.MethodDelete:
		return types.MethodDelete()
	case http.MethodConnect:
		return types.MethodConnect()
	case http.MethodOptions:
		return types.MethodOptions()
	case http.MethodTrace:
		return types.MethodTrace()
	default:
		// TODO: is other method allowed? or should we return GET?
		// https://cs.opensource.google/go/go/+/refs/tags/go1.23.2:src/net/http/method.go
		// https://github.com/WebAssembly/wasi-http/blob/main/wit/types.wit#L340C41-L341C69
		return types.MethodOther(s)
	}
}

// path returns the path with query from the [http.Request] URL.
//
// if both path and query are empty, set it to None
// [types.path]: https://github.com/WebAssembly/wasi-http/blob/main/wit/types.wit#L350
func path(req *http.Request) cm.Option[string] {
	path := req.URL.RequestURI()
	if path == "" {
		return cm.None[string]()
	}
	return cm.Some(path)
}

// toScheme returns the wasi-http [types.toScheme] from the [http.Request] URL.toScheme.
// If the scheme is not http or https, it returns a [types.SchemeOther].
//
// [types.toScheme]: https://github.com/WebAssembly/wasi-http/blob/main/wit/types.wit#L359-L360
func toScheme(s string) types.Scheme {
	switch s {
	case "http":
		return types.SchemeHTTP()
	case "https":
		return types.SchemeHTTPS()
	default:
		// TODO: when should we set the scheme to `cm.None` if `req.URL.Scheme` is empty?
		return types.SchemeOther(s)
	}
}

// toBody writes the io.ReadCloser to the wasi-http [types.toBody].
//
// [types.toBody]: https://github.com/WebAssembly/wasi-http/blob/v0.2.0/wit/types.wit#L514-L540
func toBody(body *io.ReadCloser, wasiBody *types.OutgoingBody) error {
	if body == nil || *body == nil {
		return nil
	}
	defer (*body).Close()

	stream_ := wasiBody.Write()
	stream := stream_.OK() // the first call should always return OK
	defer stream.ResourceDrop()
	if _, err := io.Copy(&outStreamWriter{stream: stream}, *body); err != nil {
		return fmt.Errorf("wasihttp: %v", err)
	}
	return nil
}

// fromBody reads the wasi-http [types.IncomingBody] to the io.ReadCloser.
//
// [types.IncomingBody]: https://github.com/WebAssembly/wasi-http/blob/v0.2.0/wit/types.wit#L397-L427
func fromBody(body *types.IncomingBody) io.ReadCloser {
	stream_ := body.Stream()
	stream := stream_.OK() // the first call should always return OK
	return &inputStreamReader{stream: stream, incomingBody: body}
}

type outStreamWriter struct {
	stream *streams.OutputStream
}

func (s *outStreamWriter) Write(p []byte) (n int, err error) {
	res := s.stream.BlockingWriteAndFlush(cm.ToList(p))
	if res.IsErr() {
		return 0, fmt.Errorf("wasihttp: %v", res.Err())
	}
	return len(p), nil
}

type inputStreamReader struct {
	stream       *streams.InputStream
	incomingBody *types.IncomingBody
}

func (r *inputStreamReader) Read(p []byte) (n int, err error) {
	poll := r.stream.Subscribe()
	poll.Block()
	poll.ResourceDrop()

	readResult := r.stream.Read(uint64(len(p)))
	if err := readResult.Err(); err != nil {
		if err.Closed() {
			return 0, io.EOF
		}
		return 0, fmt.Errorf("failed to read from InputStream %s", err.LastOperationFailed().ToDebugString())
	}

	readList := *readResult.OK()
	copy(p, readList.Slice())
	return int(readList.Len()), nil
}

func (r *inputStreamReader) Close() error {
	r.stream.ResourceDrop()
	r.incomingBody.ResourceDrop()
	return nil
}
