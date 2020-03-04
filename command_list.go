package main

import (
	"strings"
)

// CommandList is
type CommandList struct {
	commands []Command
}

// NewCommandList return new CommandList
func NewCommandList() *CommandList {
	list := new(CommandList)
	list.commands = []Command{}
	return list
}

// Add adds the command
func (list *CommandList) Add(cmd Command) {
	list.commands = append(list.commands, cmd)
}

// AppendOption append option to the command list
func (list *CommandList) AppendOption(opt string) {
	for i, cmd := range list.commands {
		list.commands[i] = append(cmd, opt)
	}
}

// Len return length of commands
func (list *CommandList) Len() int {
	return len(list.commands)
}

// String implement Stringer
func (list CommandList) String() string {
	cmds := []string{}
	for _, cmd := range list.commands {
		cmds = append(cmds, "`"+cmd.String()+"`")
	}
	return strings.Join(cmds, ", ")
}
