package discollect

import "context"

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

	return d, nil
}

// Run starts the scraping loops
func (d *Discollector) Run() error {
	w := NewWorker(d.r, d.ro, d.rl, d.q, d.w, d.er)

	w.Start()

	return nil
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
