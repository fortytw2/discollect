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

// FileWriter dumps JSON to a file
type FileWriter struct {
	f   *os.File
	enc *json.Encoder
}

// An HTTP Writer POSTS application/json to an endpoint
type HTTPWriter struct {
	c   *http.Client
	url string
}

// A MultiWriter 
type MultiWriter struct {

}