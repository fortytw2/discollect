package discollect

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// A Writer is able to process and output datums retrieved by Discollect plugins
type Writer interface {
	Write(ctx context.Context, f interface{}) error
	io.Closer
}

type StdoutWriter struct{}

func (sw *StdoutWriter) Write(ctx context.Context, f interface{}) error {
	fmt.Printf("%+v\n", f)
	return nil
}

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

// A MultiWriter writes
type MultiWriter struct {
}
