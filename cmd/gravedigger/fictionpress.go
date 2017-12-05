package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/lunny/html2md"

	"github.com/Puerkitobio/goquery"
	"github.com/fortytw2/hydrocarbon/httpx"

	dc "github.com/fortytw2/discollect"
)

// FictionPress is a plugin that can scrape fictionpress
var FictionPress = &dc.Plugin{
	Name: "fictionpress",
	ConfigValidator: func(c *dc.Config) error {
		for _, e := range c.Entrypoints {
			if !strings.Contains(e, "fictionpress.com") && !strings.Contains(e, "fanfiction.net") {
				return errors.New("fictionpress plugin only works for fictionpress and fanfiction.net")
			}
		}
		return nil
	},
	Routes: map[string]dc.Handler{
		`https:\/\/www.fictionpress.com\/s\/(.*)\/(\d+)(.*)`: storyPage,
	},
}

type chapter struct {
	Author   string    `json:"author,omitempty"`
	PostedAt time.Time `json:"posted_at,omitempty"`
	Body     string    `json:"body,omitempty"`
}

func storyPage(ctx context.Context, ho *dc.HandlerOpts, t *dc.Task) *dc.HandlerResponse {
	resp, err := ho.Client.Get(t.URL)
	if err != nil {
		return dc.ErrorResponse(err)
	}
	defer httpx.DrainAndClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return dc.ErrorResponse(errors.New("did not get 200"))
	}

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
		Body:     html2md.Convert(html.UnescapeString(strings.TrimSpace(body))),
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
				URL: fmt.Sprintf("https://www.fictionpress.com/s/%s/%s", ho.RouteParams[1], val),
			})
		})
	}

	return &dc.HandlerResponse{
		Facts: []interface{}{
			c,
		},
		Tasks: tasks,
	}
}
