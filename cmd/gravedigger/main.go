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

	err = dc.LaunchScrape("fictionpress", FictionPress.Configs[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(dc.Run())
}
