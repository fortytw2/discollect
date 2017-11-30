package main

import (
	"context"
	"fmt"
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
	fmt.Println("working task", t.URL)

	resp, err := ho.Client.Get(t.URL)
	if err != nil {
		return dc.ErrorResponse(err)
	}
	defer httpx.DrainAndClose(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return dc.ErrorResponse(err)
	}

	body, err := doc.Find(`#storytext`).Html()
	if err != nil {
		return dc.ErrorResponse(err)
	}

	c := &chapter{
		Author:   strings.TrimSpace(doc.Find(`#profile_top .xcontrast_txt+ a.xcontrast_txt`).Text()),
		PostedAt: time.Now(),
		Body:     strings.TrimSpace(body),
	}

	// find all chapters if this is the first one
	var tasks []*dc.Task
	// only for the first task
	if ho.RouteParams[2] == "1" {
		doc.Find(`#chap_select`).First().Find(`option`).Each(func(i int, sel *goquery.Selection) {
			val, exists := sel.Attr("value")
			if !exists || val == "1" {
				return
			}

			tasks = append(tasks, &dc.Task{
				URL: fmt.Sprintf("https://www.fictionpress.com/s/2961893/%s/Mother-of-Learning", val),
			})
		})
	}

	var errs []error
	return &dc.HandlerResponse{
		Errors: errs,
		Facts: []interface{}{
			c,
		},
		Tasks: tasks,
	}
}
