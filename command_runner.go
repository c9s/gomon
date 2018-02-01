package main

import (
	"os"
	"os/exec"
)

type CommandRunner struct {
	// arguments that will be passed the the commands
	args []string

	task *exec.Cmd
}

func (r *CommandRunner) IsRunning() bool {
	return r.task != nil && r.task.ProcessState != nil && !r.task.ProcessState.Exited()
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
		r.task = r.buildTask(cmd, dir)

		// Append the arguments to the task arguments
		r.task.Args = append(r.task.Args, args...)

		err := r.task.Start()
		if err != nil {
			return err
		}
		err = r.task.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CommandRunner) Stop() error {
	if r.task != nil {
		return r.task.Process.Kill()
	}
	return nil
}
