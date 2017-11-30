package discollect

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// A Worker is a single-threaded worker that pulls a single task from the queue at a time
// and process it to completion
type Worker struct {
	r  *Registry
	ro Rotator
	rl RateLimiter
	q  Queue
	w  Writer
	er ErrorReporter

	shutdown chan chan struct{}
}

// NewWorker provisions a new worker
func NewWorker(r *Registry, ro Rotator, rl RateLimiter, q Queue, w Writer, er ErrorReporter) *Worker {
	return &Worker{
		r:        r,
		ro:       ro,
		rl:       rl,
		q:        q,
		w:        w,
		er:       er,
		shutdown: make(chan chan struct{}),
	}
}

// Start launches the worker
func (w *Worker) Start() {
	for {
		select {
		case s := <-w.shutdown:
			s <- struct{}{}
			return
		default:
			qt, err := w.q.Pop(context.TODO())
			if err != nil {
				w.er.Report(context.TODO(), nil, err)
				continue
			}

			if qt == nil {
				time.Sleep(250 * time.Millisecond)
				continue
			}

			var timeout time.Duration
			if qt.Task.Timeout == timeout {
				timeout = 15 * time.Second
			} else {
				timeout = qt.Task.Timeout
			}

			// set config timeout on all worker actions on this task
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err = w.processTask(ctx, qt)
			if err != nil {
				w.er.Report(ctx, nil, err)
				cancel()
				continue
			}

			// callback that we've finished the task
			err = w.q.Finish(ctx, qt.TaskID)
			if err != nil {
				w.er.Report(ctx, nil, err)
			}

			cancel()
		}
	}
}

// Stop initiates stop and then blocks until shutdown is complete
func (w *Worker) Stop() {
	c := make(chan struct{})
	w.shutdown <- c
	<-c
}

// processTask executes one task
// Safe for concurrent use.
func (w *Worker) processTask(ctx context.Context, q *QueuedTask) error {
	handler, params, err := w.r.HandlerFor(q.Plugin, q.Task.URL)
	if err != nil {
		return err
	}

	// if this rate limit blocks too long and the context cancels we can just return
	// error and the task will be retried later
	err = w.rl.Limit(ctx, q.Config.Rate, q.Task.URL)
	if err != nil {
		return err
	}

	client, err := w.ro.Get(q.Config, q.ScrapeID)
	if err != nil {
		return err
	}

	resp := handler(ctx, &HandlerOpts{
		Config:      q.Config,
		RouteParams: params,
		Client:      client,
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
				ScrapeID: q.ScrapeID,
				Plugin:   q.Plugin,
				Config:   q.Config,
				QueuedAt: time.Now().In(time.UTC),
				Task:     t,
			}
		}

		err := w.q.Push(ctx, qt)
		if err != nil {
			errs <- err
		}
	}()

	// report errors
	go func() {
		defer wg.Done()
		for _, err := range resp.Errors {
			w.er.Report(ctx, &ReporterOpts{
				ScrapeID: q.ScrapeID,
				Plugin:   q.Plugin,
				URL:      q.Task.URL,
			}, err)
		}
	}()

	// write facts
	go func() {
		defer wg.Done()

		for _, f := range resp.Facts {
			err := w.w.Write(ctx, f)
			if err != nil {
				errs <- err
			}
		}
	}()

	// wait for all 3 writers to finish
	wg.Wait()
	// close error channel
	close(errs)

	// close error writer
	var out []error
	for e := range errs {
		if e != nil {
			out = append(out, e)
		}
	}

	if len(out) == 0 {
		return nil
	}

	return &WorkerErr{
		QueuedTask: q,
		Errors:     out,
	}
}

// WorkerErr carries errors from a task
type WorkerErr struct {
	QueuedTask *QueuedTask
	Errors     []error
}

func (we *WorkerErr) Error() string {
	return fmt.Sprintf("%v", we.Errors)
}
