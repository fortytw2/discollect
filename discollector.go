package discollect

import (
	"context"
	"errors"
)

// A Discollector ties every element of Discollect together
type Discollector struct {
	w  Writer
	r  *Registry
	rl RateLimiter
	ro Rotator
	q  Queue
	ms Metastore
	er ErrorReporter
}

// An OptionFn is used to pass options to a Discollector
type OptionFn func(d *Discollector) error

var defaultOpts = []OptionFn{
	WithWriter(&StdoutWriter{}),
	WithErrorReporter(&StdoutReporter{}),
	WithRateLimiter(&NilRateLimiter{}),
	WithRotator(&DefaultRotator{}),
	WithQueue(NewMemQueue()),
	WithMetastore(&MemMetastore{}),
}

// New returns a new Discollector
func New(opts ...OptionFn) (*Discollector, error) {
	d := &Discollector{}

	for _, o := range defaultOpts {
		err := o(d)
		if err != nil {
			return nil, err
		}
	}

	for _, o := range opts {
		err := o(d)
		if err != nil {
			return nil, err
		}
	}

	if d.r == nil {
		return nil, errors.New("no plugins registered")
	}

	return d, nil
}

// Run starts the scraping loops
func (d *Discollector) Run() error {
	w := NewWorker(d.r, d.ro, d.rl, d.q, d.w, d.er)

	w.Start()

	return nil
}

func (d *Discollector) Shutdown(ctx context.Context) {

}

// LaunchScrape starts a scrape run
func (d *Discollector) LaunchScrape(pluginName string, cfg *Config) error {
	p, err := d.r.Get(pluginName)
	if err != nil {
		return err
	}

	return launchScrape(context.TODO(), p, cfg, d.q, d.ms)
}

// WithPlugins registers a list of plugins
func WithPlugins(p ...*Plugin) OptionFn {
	return func(d *Discollector) error {
		reg, err := NewRegistry(p)
		if err != nil {
			return err
		}

		d.r = reg

		return nil
	}
}

// WithWriter sets the Writer for the Discollector
func WithWriter(w Writer) OptionFn {
	return func(d *Discollector) error {
		d.w = w
		return nil
	}
}

// WithErrorReporter sets the ErrorReporter for the Discollector
func WithErrorReporter(er ErrorReporter) OptionFn {
	return func(d *Discollector) error {
		d.er = er
		return nil
	}
}

// WithRateLimiter sets the RateLimiter for the Discollector
func WithRateLimiter(rl RateLimiter) OptionFn {
	return func(d *Discollector) error {
		d.rl = rl
		return nil
	}
}

// WithRotator sets the Rotator for the Discollector
func WithRotator(ro Rotator) OptionFn {
	return func(d *Discollector) error {
		d.ro = ro
		return nil
	}
}

// WithQueue sets the Queue for the Discollector
func WithQueue(q Queue) OptionFn {
	return func(d *Discollector) error {
		d.q = q
		return nil
	}
}

// WithMetastore sets the Metastore for the Discollector
func WithMetastore(ms Metastore) OptionFn {
	return func(d *Discollector) error {
		d.ms = ms
		return nil
	}
}
