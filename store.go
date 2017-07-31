package discollect

import "context"

type Store interface {
	// StartScrape attempts to start the scrape, returning `true, nil` if the scrape is
	// able to be started
	StartScrape(ctx context.Context, pluginName string, cfg *Config) (id string, ok bool, err error)
	EndScrape(ctx context.Context, id string, facts, tasks int) error
}
