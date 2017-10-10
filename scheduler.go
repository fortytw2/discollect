package discollect

// A Scheduler initiates new scrapes according to plugin-level schedules
type Scheduler struct {
	r  *Registry
	ms Metastore
	q  Queue
}

// Start launches the scheduler
func (s *Scheduler) Start() {

}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {

}

// A Schedule is part of every plugin and defines when it needs to be run
type Schedule struct {
	Cron string
}
