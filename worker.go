package discollect

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// A Worker is a single-threaded worker that pulls N tasks from the queue and processes them
// in order. To achieve parellelism, use a WorkerGroup
type Worker struct {
	r       *Registry
	rotator Rotator
	q       Queue
	writer  Writer
	er      ErrorReporter

	shutdown chan struct{}
	closed   chan struct{}
}

// NewWorker provisions a new worker
func NewWorker(r *Registry, ro Rotator, q Queue, w Writer, er ErrorReporter) *Worker {
	return &Worker{
		shutdown: make(chan struct{}),
		closed:   make(chan struct{}),
	}
}

// Start launches the worker
func (w *Worker) Start(pullCount int) {
	defer func() {
		w.closed <- struct{}{}
	}()

	for {
		select {
		case <-w.shutdown:
			return
		default:
			qts, err := w.q.Pop(context.TODO(), pullCount)
			if err != nil {
				w.er.Report(context.TODO(), nil, err)
				continue
			}

			for _, qt := range qts {
				err = w.processTask(context.TODO(), qt)
				if err != nil {
					w.er.Report(context.TODO(), nil, err)
					continue
				}

				// callback that we've finished the task
				err = w.q.MarkDone(context.TODO(), qt.TaskID)
				if err != nil {
					w.er.Report(context.TODO(), nil, err)
				}
			}
		}
	}
}

// Stop initiates stop and then blocks until shutdown is complete
func (w *Worker) Stop() {
	w.shutdown <- struct{}{}
	<-w.closed
}

// processTask executes one task
// Safe for concurrent use.
func (w *Worker) processTask(ctx context.Context, q *QueuedTask) error {
	handler, err := w.r.HandlerFor(q.Plugin, q.Task.URL)
	if err != nil {
		return err
	}

	client, err := w.rotator.Get(q.Config, q.ScrapeID)
	if err != nil {
		return err
	}

	resp := handler(ctx, &HandlerOpts{
		Config: q.Config,
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
			err := w.writer.Write(ctx, f)
			if err != nil {
				errs <- err
			}
		}
	}()

	// wait for all 3 writers to finish
	wg.Done()

	// close error writer
	var out []error
	for e := range errs {
		out = append(out, e)
	}

	if len(out) == 0 {
		return nil
	}

	return &WorkerErr{
		QueuedTask: q,
		Errors:     out,
	}
}

type WorkerErr struct {
	QueuedTask *QueuedTask
	Errors     []error
}

func (we *WorkerErr) Error() string {
	return fmt.Sprintf("%v", we.Errors)
}
