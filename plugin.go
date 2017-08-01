package discollect

import (
	"context"
	"net/http"
	"time"

	"github.com/fortytw2/discollect/countries"
)

type Config struct {
	// DynamicLaunch specifies whether this plugin must have preloaded
	// LaunchInfo, or if they can be dynamically specified, i.e by an end user
	DynamicLaunch bool
	// LaunchInfo is used to start a scrape
	LaunchInfo []string
	// A Plugin can have many types of Scrapes
	// common ones are "standard" and "quick"
	Type string
	// Since is used to convey delta information
	Since time.Time
	// Rate is used to configure rate limits, per-scrape, per-ip, and per-domain
	Rate *RateLimit
	// Countries is a list of countries this scrape can be executed from
	Countries []countries.Country
}

// RateLimit is a wrapper struct around a variety of per-config rate limits
type RateLimit struct {
	PerIP     int
	PerScrape int
	PerDomain int
}

type Plugin struct {
	Name     string
	Schedule *Schedule
	Configs  []*Config
	Routes   map[string]Handler
}

// HandlerOpts are passed to
type HandlerOpts struct {
	Config *Config
	Client *http.Client
}

type HandlerResponse struct {
	Tasks  []*Task
	Facts  []Fact
	Errors []error
}

type Handler func(ctx context.Context, ho *HandlerOpts, t *Task) *HandlerResponse
