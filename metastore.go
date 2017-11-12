package discollect

import "context"

// A Metastore is used to store the history of all scrape runs
type Metastore interface {
	// StartScrape attempts to start the scrape, returning `true, nil` if the scrape is
	// able to be started
	StartScrape(ctx context.Context, pluginName string, cfg *Config) (id string, err error)
	EndScrape(ctx context.Context, id string, datums, tasks int) error
}
