package main

import (
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

func (r *CommandRunner) buildTask(cmd Command, dir string) *exec.Cmd {
	p := exec.Command(cmd[0], cmd[1:]...)
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	p.Dir = dir
	return p
}

func (r *CommandRunner) Run(commands []Command, args []string, dir string) error {
	for _, cmd := range commands {
		{
			var task = r.buildTask(cmd, dir)
			task.Args = append(task.Args, args...)

			r.mu.Lock()
			r.task = task
			r.mu.Unlock()

			if err := task.Start(); err != nil {
				return err
			}

			if err := task.Wait(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *CommandRunner) Stop() error {
	var task = r.Task()
	if task != nil {
		return task.Process.Kill()
	}
	return nil
}
