package main

import (
	"context"
	"log"
	"net/http"

	"github.com/fortytw2/discollect"
	"github.com/fortytw2/discollect/api"
	"github.com/oklog/oklog/pkg/group"
)

func main() {
	dc, err := discollect.New(
		discollect.WithPlugins(FictionPress),
	)
	if err != nil {
		log.Fatal(err)
	}

	// err = dc.LaunchScrape("fictionpress", &discollect.Config{
	// 	DynamicEntry: true,
	// 	Entrypoints:  []string{`https://www.fictionpress.com/s/2961893/1/Mother-of-Learning`},
	// 	Type:         "full",
	// 	Name:         "Mother of Learning",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	r := api.Router(dc)
	h := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	var g group.Group
	{
		g.Add(h.ListenAndServe, func(error) {
			h.Shutdown(context.Background())
		})
	}
	{
		g.Add(func() error { return dc.Start(5) }, func(error) {
			dc.Shutdown(context.Background())
		})
	}

	log.Fatal(g.Run())
}
