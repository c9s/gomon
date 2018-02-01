package main

import (
	"strings"
)

type Command []string

func (cmd *Command) String() string {
	return strings.Join(*cmd, " ")
}
