package discollect

import "context"

// A ProgressTracker is able to track realtime progress of a given scrape
type ProgressTracker interface {
	AddTasks(ctx context.Context, scrapeID string, tasks int) error
	FinishTasks(ctx context.Context, scrapeID string, tasks int) error
}
