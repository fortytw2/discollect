package discollect

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

// A Queue is used to submit and retrieve individual tasks
type Queue interface {
	Pop(ctx context.Context) (*QueuedTask, error)
	Push(ctx context.Context, tasks []*QueuedTask) error

	Finish(ctx context.Context, taskID ulid.ULID) error
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

// NewMemQueue makes a new purely in-memory queue
func NewMemQueue() *MemQueue {
	return &MemQueue{}
}

type MemQueue struct {
	mu sync.Mutex
	q  []*QueuedTask
}

func (mq *MemQueue) Pop(ctx context.Context) (*QueuedTask, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if len(mq.q) == 0 {
		return nil, nil
	}

	qt := mq.q[0]
	mq.q = mq.q[1:]

	return qt, nil
}
func (mq *MemQueue) Push(ctx context.Context, tasks []*QueuedTask) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.q = append(mq.q, tasks...)
	return nil
}

func (mq *MemQueue) Finish(ctx context.Context, taskID ulid.ULID) error {
	return nil
}
