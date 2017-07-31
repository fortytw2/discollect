package discollect

import "net/http"

// Rotator is a proxy rotator interface idk
type Rotator interface {
	Get(c *Config, id string) (*http.Client, error)
}
