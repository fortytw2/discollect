package discollect

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	r       *Registry
	rotator Rotator
	q       Queue
	writer  Writer
	er      ErrorReporter
}

// Work executes one task
// Safe for concurrent use.
func (w *Worker) Work(ctx context.Context, id string, cfg *Config, q *QueuedTask) error {
	handler, err := w.r.HandlerFor(q.Plugin, q.Task.URL)
	if err != nil {
		return err
	}

	client, err := w.rotator.Get(cfg, id)
	if err != nil {
		return err
	}

	resp := handler(ctx, &HandlerOpts{
		Config: cfg,
		Client: client,
	}, q.Task)

	errs := make(chan error, 8)
	var wg sync.WaitGroup
	wg.Add(3)
	// push queued tasks
	go func() {
		defer wg.Done()

		qt := make([]*QueuedTask, len(resp.Tasks))
		for i, t := range resp.Tasks {
			qt[i] = &QueuedTask{
				ScrapeID: id,
				Plugin:   q.Plugin,
				QueuedAt: time.Now().In(time.UTC),
				Task:     t,
			}
		}

		err := w.q.PushN(ctx, qt)
		if err != nil {
			errs <- err
		}
	}()

	// report errors
	go func() {
		defer wg.Done()
		for _, err := range resp.Errors {
			w.er.Report(ctx, &ReporterOpts{
				ScrapeID: id,
				Plugin:   q.Plugin,
			}, err)
		}
	}()

	// write facts
	go func() {
		defer wg.Done()

		for _, f := range resp.Facts {
			err := w.writer.Write(ctx, f)
			if err != nil {
				errs <- err
			}
		}
	}()

	var out []error
	go func() {
		for e := range errs {
			out = append(out, e)
		}
	}()

	// wait for all 3 writers to finish
	wg.Done()
	// close error writer
	close(errs)

	if len(out) == 0 {
		return nil
	}

	return &WorkerErr{
		Errors: out,
	}
}

type WorkerErr struct {
	Errors []error
}

func (we *WorkerErr) Error() string {
	return fmt.Sprintf("%v", we.Errors)
}
