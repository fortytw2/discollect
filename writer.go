package discollect

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// A Writer is able to process and output datums retrieved by Discollect plugins
type Writer interface {
	Write(ctx context.Context, f interface{}) error
	io.Closer
}

// StdoutWriter fmt.Printfs to stdout
type StdoutWriter struct{}

// Write printf %+v the datum to stdout
func (sw *StdoutWriter) Write(ctx context.Context, f interface{}) error {
	return json.NewEncoder(os.Stdout).Encode(f)
}

// Close is a no-op function so the StdoutWriter works
func (sw *StdoutWriter) Close() error {
	return nil
}

// FileWriter dumps JSON to a file
type FileWriter struct {
	f   *os.File
	enc *json.Encoder
}

// A HTTPWriter POSTS application/json to an endpoint
type HTTPWriter struct {
	c   *http.Client
	url string
}

// A MultiWriter writes to multiple writers at once, in parallel
type MultiWriter struct {
}
