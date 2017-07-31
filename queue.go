package discollect

import (
	"context"
	"time"
)

type Queue interface {
	PopN(ctx context.Context, i int) ([]*QueuedTask, error)
	PushN(ctx context.Context, tasks []*QueuedTask) error
}

type QueuedTask struct {
	ScrapeID string    `json:"scrape_id"`
	QueuedAt time.Time `json:"queued_at"`
	Plugin   string    `json:"plugin"`
	Task     *Task     `json:"task"`
}

// A Task generally maps to a single HTTP request, but sometimes more than one
// may be made
type Task struct {
	URL   string            `json:"url"`
	Extra map[string]string `json:"extra"`
}
