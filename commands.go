package main

var goCommands = map[string][]string{
	"test":    []string{"go", "test"},
	"install": []string{"go", "install"},
	"build":   []string{"go", "build"},
	"fmt":     []string{"go", "fmt"},
	"fix":     []string{"go", "fix"},
	"run":     []string{"go", "run"},
	"vet":     []string{"go", "vet"},
}
