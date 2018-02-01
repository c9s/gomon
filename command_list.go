package main

import (
	"os/exec"
	"strings"
)

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
		list.commands[i] = append(cmd, opt)
	}
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
