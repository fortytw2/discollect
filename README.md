discollect
------

dis' collector. 

distributed web scraper inspired by scrapy

Architecture
------

`discollect` is designed to power known-entrypoint, structured scrapes,
unlike a "web crawler". You must have an entrypoint - an initial URL(s) 
that you're going to submit as a `Task`. From there, every URL you submit 
is able to be matched to a `Handler`, much like you would match routes
in a web server and each handler returns more `Task`s and datums, or an error.

By changing the entrypoints for a given plugin, we're able to change the nature
of a scrape - collecting small amounts of priority information from a site using
the same code that drives a full site scrape. i.e. for themoviedb.org, instead of
entering at the the published list of all movie IDs on the site, we could enter
at an individual movie level.

This architecture allows for large scale, distributed, structured scraping of a 
variety of websites.

These tasks are stored in the `Task Queue`, while information and statistics about
individual jobs are stored in the `Metastore`. Also in the `Metastore` are detailed
error reports from when `Handler`s return errors during execution and fail to retry.

Getting Started (Library Use)
------

`discollect` is fundamentally a very flexible library - 99% of the time users will 
want to implement their own plugins, writers, metastores and task queues to better
fit into their environment. Natively, `discollect` comes with many building blocks 
to make this as easy as possible. A basic setup may be as easy as 

```go
package main

import (
    "log"

    "github.com/fortytw2/discollect"
    "github.com/fortytw2/discollect/metastore/sqlite"
    "github.com/your-org/plugins"
)

func main() {
    d, err := discollect.New(
        discollect.WithWriter(discollect.NewFileWriter("/output")),
        discollect.WithTaskQueue(discollect.NewInMemTaskQueue()),
        discollect.WithMetastore(sqlite.NewMetastore("discollect.db")),
        discollect.WithErrorReporter(discollect.)
    )

    err := d.RegisterPlugins(
        plugins.NewTMDBPlugin("api-key-here"),
        plugins.NewBlogCrawler(),
    )

    // parse all CLI options
    d.ParseFlags()

    // handle SIGTERM/SIGKILL gracefully
    go d.SignalHandler()

    // serve HTTP
    log.Fatal(d.ListenAndServe())
}
```

However, `discollect` also ships with a binary `discollectd`, which includes all task 
queues, metastores, and writers that are part of the core library. If you build 
plugins entirely in `skylark` or `lua` and do not need to implement any custom queues,
metastores, or writers, this may be a good route to go.

Plugins
------

Plugins are written in either `go`, `skylark`, or `lua`. Go plugins require a recompilation
of `discollect` to be deployed, whereas plugins built in `skylark` or `lua` simply require
a reload.

`discollect` includes a large number of test helper functions to aid you in writing effective,
fast tests for your network-based plugins. Read the docs to learn more. (DOC LINK HERE)

Metastores
------

A `metastore` is the key to a `discollect` deployment, collecting task and scrape information
as well as handling leader election (in a distributed, HA environment)

As always, `discollect` ships with a variety of `metastore` implementations for you to choose 
from

- SQLite3 (good for single node deployments, can work for worker deployments)
- PostgreSQL 9.4+ (for single master, many worker deployments)
- etcd3 (for multi-master, true HA deployments)

Task Queues
------

A reliable task queue is the heart of any distributed work system. For `discollect`, there 
are many options to choose from

- Disk-backed, durable queue (only for single node deployments)
- Redis 3+ (BRPOPLPUSH based)
- Beanstalkd
- Amazon SQS 
- Google Cloud Pub/Sub

Writers
------

When a datum is returned from a `Handler`, `discollect` can perform several actions natively 

- POST an endpoint over HTTP with the datum encoded as JSON
- write the datum to a file per scrape run

However, it is very easy to extend `discollect` to perform any action with a datum by implementing 
the simple `Writer` interface - 

```go
// A Writer is able to process and output datums retrieved by Discollect plugins
type Writer interface {
	Write(ctx context.Context, datum interface{}) error
	io.Closer
}
```

Writers can also be composed using the `discollect.NewMultiWriter(w ...Writer) Writer` 
helper function.

Scheduled Scrapes
------

`discollect` packs a powerful, reliable, cron-based scheduler into every deployment
that tightly integrates with plugin authorship and the built-in alerting capabilites.

Alerting and Monitoring
------

`discollect` is able to monitor internal scrape stability (i.e. number of tasks and datums
per run of the same plugin and config over time) and alert to several different providers

- Slack
- IRC
- Pagerduty

Region-Aware Plugins & Proxy Rotation
------

For many scraping tasks, it's neccessary to only make requests from certain IPs.


License
------
