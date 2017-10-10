package discollect

// A Discollector ties every element of Discollect together
type Discollector struct {
	w  Writer
	r  *Registry
	ro Rotator
	q  Queue
	ms Metastore
	er ErrorReporter
}

// An OptionFn is used to pass options to a Discollector
type OptionFn func(d *Discollector) error

// New returns a new Discollector
func New(opts ...OptionFn) (*Discollector, error) {
	return nil, nil
}
