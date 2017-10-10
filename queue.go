package discollect

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fortytw2/discollect/testing"
)

// A Queue is used to submit and retrieve individual tasks
type Queue interface {
	Pop(ctx context.Context, i int) ([]*QueuedTask, error)
	Push(ctx context.Context, tasks []*QueuedTask) error

	MarkDone(ctx context.Context, taskID string) error
}

type QueuedTask struct {
	// set by the TaskQueue
	TaskID   string `json:"task_id"`
	ScrapeID string `json:"scrape_id"`

	QueuedAt time.Time `json:"queued_at"`
	Config   *Config   `json:"config"`
	Plugin   string    `json:"plugin"`
	Retries  int       `json:"retries"`

	Task *Task `json:"task"`
}

// A Task generally maps to a single HTTP request, but sometimes more than one
// may be made
type Task struct {
	URL string `json:"url"`
	// Extra can be used to send information from a parent task to its children
	Extra map[string]json.RawMessage `json:"extra"`
}

// TestQueue is a re-useable conformance test for any implementation of Queue to pass
func TestQueue(t testing.T, q Queue) {
	return
}
