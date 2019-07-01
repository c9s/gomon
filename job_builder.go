package main

import (
	"path/filepath"
)

type JobBuilder struct {
	// Job template arguments
	Commands []Command
	Args     []string

	// template options
	ChangeDirectory bool
	AppendFilename  bool
}

func (t *JobBuilder) Create(filename string) *Job {
	var chdir = ""
	if t.ChangeDirectory {
		chdir = filepath.Dir(filename)
	}

	var args []string

	copy(args, t.Args)

	if t.AppendFilename {
		args = append(args, filename)
	}

	return &Job{
		commands: t.Commands,
		args:     args,
		dir:      chdir,
	}
}
