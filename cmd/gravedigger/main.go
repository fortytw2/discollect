package main

import (
	"log"

	"github.com/fortytw2/discollect"
)

func main() {
	dc, err := discollect.New(
		discollect.WithPlugins(FictionPress),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = dc.LaunchScrape("fictionpress", &discollect.Config{
		DynamicEntry: true,
		Entrypoints:  []string{`https://www.fictionpress.com/s/2961893/1/Mother-of-Learning`},
		Type:         "full",
		Name:         "Mother of Learning",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(dc.Run())
}
