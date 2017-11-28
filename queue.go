package discollect

import (
	"context"
	"encoding/json"
	"time"

	"github.com/oklog/ulid"
)

// A Queue is used to submit and retrieve individual tasks
type Queue interface {
	Pop(ctx context.Context) (*QueuedTask, error)
	Push(ctx context.Context, tasks []*QueuedTask) error

	Finish(ctx context.Context, taskID ulid.ULID) error
	Retry(context.Context, *QueuedTask) error
}

// A QueuedTask is the struct for a task that goes on the Queue
type QueuedTask struct {
	// set by the TaskQueue
	TaskID   ulid.ULID `json:"task_id"`
	ScrapeID string    `json:"scrape_id"`

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
	Extra map[string]json.RawMessage `json:"extra,omitempty"`
	// Timeout is the timeout a single task should have attached to it
	// defaults to 15s
	Timeout time.Duration
}
