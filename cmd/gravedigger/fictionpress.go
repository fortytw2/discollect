package main

import (
	"context"
	"strings"
	"time"

	"github.com/Puerkitobio/goquery"
	"github.com/fortytw2/hydrocarbon/httpx"

	dc "github.com/fortytw2/discollect"
)

// FictionPress is a plugin that can scrape fictionpress
var FictionPress = &dc.Plugin{
	Name: "fictionpress",
	Configs: []*dc.Config{
		{
			Entrypoints: []string{`https://www.fictionpress.com/s/2961893/1/Mother-of-Learning`},
			Type:        "full",
		},
	},
	Routes: map[string]dc.Handler{
		`https:\/\/www.fictionpress.com\/s\/(.*)\/(\d+)\/(.*)`: storyPage,
	},
}

type chapter struct {
	Author   string
	PostedAt time.Time
	Body     string
}

func storyPage(ctx context.Context, ho *dc.HandlerOpts, t *dc.Task) *dc.HandlerResponse {
	resp, err := ho.Client.Get(t.URL)
	if err != nil {
		return dc.ErrorResponse(err)
	}
	defer httpx.DrainAndClose(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return dc.ErrorResponse(err)
	}

	c := &chapter{
		Author:   strings.TrimSpace(doc.Find(`#profile_top .xcontrast_txt+ a.xcontrast_txt`).Text()),
		PostedAt: time.Now(),
		Body:     strings.TrimSpace(doc.Find(`#storytext`).Text()),
	}

	var errs []error
	return &dc.HandlerResponse{
		Errors: errs,
		Facts: []interface{}{
			c,
		},
		// Tasks: []*dc.Task{
		// 	{
		// 		URL: "",
		// 	},
		// },
	}
}
