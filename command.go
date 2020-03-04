package main

import (
	"strings"
)

// Command is
type Command []string

// String implements Stringer
func (cmd *Command) String() string {
	return strings.Join(*cmd, " ")
}
