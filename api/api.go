package api

import (
	"net/http"

	"github.com/fortytw2/discollect"
)

// Router returns an *http.ServeMux set up to expose the functionality of the discollector over HTTP
func Router(dc *discollect.Discollector) *http.ServeMux {
	m := http.NewServeMux()

	m.Handle("/scrapes", scrapeHandler(dc))

	return m
}

func scrapeHandler(dc *discollect.Discollector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listScrapes(dc, w, r)
		case http.MethodPost:
			startScrape(dc, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func listScrapes(dc *discollect.Discollector, w http.ResponseWriter, r *http.Request) {

}

func startScrape(dc *discollect.Discollector, w http.ResponseWriter, r *http.Request) {

}
