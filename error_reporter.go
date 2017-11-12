package discollect

import (
	"context"
)

// An ErrorReporter is used to send forward faulty handler runs to
// a semi-permanent sink for later analysis.
// Generally, this can be a service such as Sentry or Bugsnag
// but may also be a simpler DB backend, like Postgres
// An ErrorReporter should discard any calls with err == nil
type ErrorReporter interface {
	Report(ctx context.Context, ro *ReporterOpts, err error)
}

// ReporterOpts is used to attach additional information to an error
type ReporterOpts struct {
	ScrapeID string
	Plugin   string
	URL      string
}
