package main

import (
	"context"
	"os"
	"os/exec"
	"sync"
)

type CommandRunner struct {
	// args is an argument array that will be passed the the commands
	args []string

	// task is the current executing command
	task *exec.Cmd

	mu sync.Mutex
}

func (r *CommandRunner) Task() (task *exec.Cmd) {
	r.mu.Lock()
	task = r.task
	r.mu.Unlock()
	return task
}

func (r *CommandRunner) IsRunning() bool {
	var task = r.Task()
	return task != nil && task.ProcessState != nil && !task.ProcessState.Exited()
}

func buildTask(ctx context.Context, cmd Command, dir string, args []string) *exec.Cmd {
	var p = exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	p.Dir = dir
	p.Args = append(p.Args, args...)
	return p
}

func (r *CommandRunner) Run(ctx context.Context, commands []Command, args []string, dir string) error {
	for _, cmd := range commands {
		var task = buildTask(ctx, cmd, dir, args)

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

func (r *CommandRunner) Stop() (err error) {
	r.mu.Lock()
	if r.task != nil {
		err = r.task.Process.Kill()
		r.task = nil
	}
	r.mu.Unlock()
	return err
}
