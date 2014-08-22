package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var goCommands = map[string][]string{
	"test":    []string{"go", "test"},
	"install": []string{"go", "install"},
	"build":   []string{"go", "build"},
	"fmt":     []string{"go", "fmt"},
	"fix":     []string{"go", "fix"},
	"run":     []string{"go", "run"},
	"vet":     []string{"go", "vet"},
}

type Command []string

func (cmd *Command) String() string {
	return strings.Join(*cmd, " ")
}

func (cmd *Command) buildTask(dir *string) *exec.Cmd {
	task := exec.Command((*cmd)[0], (*cmd)[1:]...)
	task.Stdout = os.Stdout
	task.Stderr = os.Stderr
	if dir != nil {
		task.Dir = *dir
	}
	return task
}

type CommandList struct {
	commands []Command
	task     *exec.Cmd
	filename string
}

func NewCommandList() *CommandList {
	list := new(CommandList)
	list.commands = []Command{}
	return list
}

func (list *CommandList) Add(cmd Command) {
	list.commands = append(list.commands, cmd)
}

func (list *CommandList) AppendOption(opt string) {
	for i, cmd := range list.commands {
		list.commands[i] = append(cmd, "-x")
	}
}

func (list *CommandList) ClearFilename() {
	list.filename = ""
}

func (list *CommandList) SetFilename(filename string) {
	list.filename = filename
}

func (list *CommandList) Len() int {
	return len(list.commands)
}

func (list CommandList) String() string {
	cmds := []string{}
	for _, cmd := range list.commands {
		cmds = append(cmds, "`"+cmd.String()+"`")
	}
	return strings.Join(cmds, ", ")
}

func (list *CommandList) IsTaskRunning() bool {
	return list.task != nil && list.task.ProcessState != nil && !list.task.ProcessState.Exited()
}

func (list *CommandList) StopTask() error {
	if list.IsTaskRunning() {
		fmt.Println("Stopping Task...")
		return list.task.Process.Kill()
	}
	return nil
}

func (list *CommandList) Run(dir *string) error {
	if list.IsTaskRunning() {
		return errors.New("Previous task is still running")
	}
	for _, cmd := range list.commands {
		list.task = cmd.buildTask(dir)
		if list.filename != "" {
			list.task.Args = append(list.task.Args, list.filename)
		}
		err := list.task.Start()
		if err != nil {
			return err
		}
		err = list.task.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}
