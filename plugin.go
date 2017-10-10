package discollect

import (
	"context"
	"net/http"
	"time"

	"github.com/fortytw2/discollect/countries"
)

type Config struct {
	// DynamicEntry specifies whether this plugin must have preloaded
	// Entrypoints, or if they can be dynamically specified, i.e by an end user
	DynamicEntry bool
	// Entrypoints is used to start a scrape
	Entrypoints []string
	// A Plugin can have many types of Scrapes
	// common ones are "full" and "delta"
	Type string
	// Since is used to convey delta information
	Since time.Time
	// Rate is used to configure rate limits, per-scrape, per-ip, and per-domain
	Rate *RateLimit
	// Countries is a list of countries this scrape can be executed from
	// nil if unused
	Countries []countries.Country
}

// RateLimit is a wrapper struct around a variety of per-config rate limits
type RateLimit struct {
	// Rate a single IP can make requests
	PerIP int
	// Rate the entire scrape can operate at
	PerScrape int
	// Rate per domain using the publicsuffix list to differentiate
	PerDomain int
}

// A Plugin is capable of running scrapes, ideally of a common type or against a single site
type Plugin struct {
	Name     string
	Schedule *Schedule
	Configs  []*Config
	Routes   map[string]Handler
}

// HandlerOpts are passed to a Handler
type HandlerOpts struct {
	Config *Config
	Client *http.Client
}

// HandlerResponse
type HandlerResponse struct {
	Tasks  []*Task
	Facts  []interface{}
	Errors []error
}

type Handler func(ctx context.Context, ho *HandlerOpts, t *Task) *HandlerResponse
