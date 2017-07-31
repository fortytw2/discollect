package discollect

import "context"

type Fact interface{}

// A Writer is able to process and output Facts retrieved by Discollect plugins
type Writer interface {
	Write(ctx context.Context, f Fact) error
}
