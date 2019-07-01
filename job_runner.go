package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/c9s/gomon/logger"
)

type JobRunner struct {
	builder *JobBuilder

	lastJob          *Job
	previouslyFailed bool

	mu     sync.Mutex
	ctx    context.Context
	cancel func()
}

func (r *JobRunner) RunAndNotify(ctx context.Context, filename string, alwaysNotify bool) (duration time.Duration, err error) {
	duration, err = r.Run(ctx, filename)

	if err != nil {
		r.mu.Lock()
		r.previouslyFailed = true
		r.mu.Unlock()

		notifier.NotifyFailed("Build failed", err.Error())
	} else {
		r.mu.Lock()
		if r.previouslyFailed {
			r.previouslyFailed = false

			notifier.NotifyFixed("Build fixed", fmt.Sprintf("Spent: %s", duration))
		} else if alwaysNotify {
			notifier.NotifySucceeded("Build succeeded", fmt.Sprintf("Spent: %s", duration))
		}
		r.mu.Unlock()
	}
	return
}

func (r *JobRunner) Run(basectx context.Context, filename string) (duration time.Duration, err error) {
	r.mu.Lock()

	if r.ctx != nil {
		logger.Warnln("Canceling previous context")
		r.cancel()
		r.ctx = nil
		r.cancel = nil
	}
	if r.lastJob != nil {
		logger.Infof("Stopping job: %v", r.lastJob)
		if err := r.lastJob.StopAndWait(); err != nil {
			logger.Errorf("Failed to stop job. error=%v", err)
		}
	}

	// allocate a new context
	r.ctx, r.cancel = context.WithCancel(basectx)

	var ctx = r.ctx
	var job = r.builder.Create(filename)

	r.lastJob = job
	r.mu.Unlock()

	logger.Infof("Starting: commands=%v args=%v", job.commands, job.args)
	var now = time.Now()
	err = job.Run(ctx)
	duration = time.Now().Sub(now)

	r.mu.Lock()
	r.lastJob = nil
	r.ctx = nil
	r.cancel = nil
	r.mu.Unlock()

	if err != nil {
		return duration, err
	}

	return duration, nil
}
