package api

import (
	"github.com/fortytw2/discollect"
)

// A Server is able to serve an HTTP/JSON API on top of a discollector
type Server struct {
	d *discollect.Discollector
}
