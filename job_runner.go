package main

import (
	"context"
	"sync"
	"time"

	"github.com/c9s/gomon/logger"
)

type JobRunner struct {
	builder *JobBuilder

	lastJob *Job

	mu     sync.Mutex
	ctx    context.Context
	cancel func()
}

func (r *JobRunner) Run(filename string) (duration time.Duration, err error) {
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
	r.ctx, r.cancel = context.WithCancel(context.Background())

	var ctx = r.ctx
	var job = r.builder.Create(filename)

	r.lastJob = job
	r.mu.Unlock()

	logger.Infof("Starting: commands=%v args=%v", job.commands, job.args)
	var now = time.Now()
	err = job.Run(ctx)
	duration = time.Now().Sub(now)
	if err != nil {
		return duration, err
	}

	return duration, nil
}
