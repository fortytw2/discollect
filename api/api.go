package api

import (
	"net/http"

	"github.com/fortytw2/discollect"
)

func Router(dc *discollect.Discollector) *http.ServeMux {
	return http.DefaultServeMux
}
