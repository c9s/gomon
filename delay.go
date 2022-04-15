package main

import (
	"context"
	"sync"
	"time"

	"github.com/c9s/gomon/logger"
)

const delayWait time.Duration = 500 * time.Millisecond

type DelayResult struct {
	duration time.Duration
	err      error
}

// Wraps JobRunner so that we only call it every delayWait
// duration. Used since file changes often come in batches
// (e.g. go fmt & the vim plugin), and since we also want to support
// cancelling running jobs.
type Delay struct {
	jobRunner    *JobRunner
	alwaysNotify bool
	jobs         chan string
	results      chan DelayResult
	wait         sync.WaitGroup
}

func NewDelay(jobRunner *JobRunner, alwaysNotify bool) *Delay {
	return &Delay{
		jobRunner:    jobRunner,
		alwaysNotify: alwaysNotify,
		jobs:         make(chan string),
		results:      make(chan DelayResult),
	}
}

// Logs results of command run after completion, with
// a little extra logic to detect cancellation.
func (d *Delay) logResult(ctx context.Context) {
	var result DelayResult
	var interrupted bool
	select {
	case <-ctx.Done():
		// Context was cancelled before result came in.
		// We must still wait for results, though!
		interrupted = true
		result = <-d.results
	case result = <-d.results:
	}

	switch {
	case interrupted:
		logger.Infoln("Command was killed and restarted...")
	case result.err != nil:
		logger.Errorf("Command failed: %v", result.err.Error())
	default:
		logger.Infoln("Command succeeded:", result.duration)
	}
}

// Triggers delayed execution for a given file event.
// NB: only the most recent file in a 500ms batch will be used.
func (d *Delay) Trigger(filename string) {
	d.jobs <- filename
}

// Rolls up triggers into a single run of the job runner every
// delayWait duration, with logic for cancelling any active job.
func (d *Delay) Run() {
	filename := ""
	ctx, cancel := context.WithCancel(context.Background())
	for {
		select {
		case <-time.After(delayWait):
			if filename == "" {
				continue
			}

			//Cancel any running job and wait for it to exit.
			cancel()
			d.wait.Wait()

			// Construct new context and execute job.
			ctx, cancel = context.WithCancel(context.Background())
			go func(ctx context.Context, filename string) {
				d.wait.Add(1)
				defer cancel()
				defer d.wait.Done()
				duration, err := d.jobRunner.RunAndNotify(ctx, filename, d.alwaysNotify)
				d.results <- DelayResult{duration, err}
			}(ctx, filename)

			go d.logResult(ctx)
			filename = ""

		case filename = <-d.jobs:
		}
	}
}
