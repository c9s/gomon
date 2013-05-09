package main

var goCommands = map[string][]string{
	"test":    []string{"go", "test"},
	"install": []string{"go", "install"},
	"build":   []string{"go", "build"},
	"fmt":     []string{"go", "fmt"},
	"run":     []string{"go", "run"},
}

type Command []string

type CommandSet struct {
	Commands []Command
}
