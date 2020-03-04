package main

import (
	"context"
	"os"
	"os/exec"
	"sync"
)

// Job is
type Job struct {
	commands []Command
	args     []string
	dir      string

	// currentTask is the current executing command
	currentTask *exec.Cmd
	mu          sync.Mutex
}

// IsRunning return whether the process still alive
func (job *Job) IsRunning() bool {
	var task = job.currentTask
	return task != nil && task.ProcessState != nil && !task.ProcessState.Exited()
}

// Run spawn the command with context
func (job *Job) Run(ctx context.Context) error {
	for _, cmd := range job.commands {
		job.mu.Lock()
		var task = buildTask(ctx, cmd, job.args)
		task.Dir = job.dir
		job.currentTask = task
		job.mu.Unlock()

		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := task.Start(); err != nil {
			return err
		}

		if err := task.Wait(); err != nil {
			return err
		}
	}
	return nil
}

// StopAndWait request to stop the Job and want the process exited
func (job *Job) StopAndWait() (err error) {
	job.mu.Lock()
	defer job.mu.Unlock()

	if job.currentTask != nil {
		task := job.currentTask

		if !task.ProcessState.Exited() {
			err = task.Process.Kill()
			if err != nil {
				return err
			}

			err = task.Wait()
			if err != nil {
				return err
			}
		}

		job.currentTask = nil
	}
	return err
}

func buildTask(ctx context.Context, cmd Command, args []string) *exec.Cmd {
	var p = exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	p.Args = append(p.Args, args...)
	return p
}
