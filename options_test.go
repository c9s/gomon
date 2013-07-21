package main

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	opt := Option{"h", false, "Show Help"}
	opts := Options{&opt}
	if opts.Get("h") != &opt {
		t.Fatal("Get failed")
	}
	if opts.Get("unknown") != nil {
		t.Fatal("Get should return nil for unknown option")
	}
}

func TestParseNoArg(t *testing.T) {
	args := strings.Split("gomon", " ")
	opts := NewOptions()
	dirArgs, cmdArgs := opts.Parse(args)
	if len(dirArgs) != 0 {
		t.Fatal("dirArgs should be empty")
	}
	if len(cmdArgs) != 0 {
		t.Fatal("cmdArgs should be empty")
	}
}

func TestParseDirArg(t *testing.T) {
	args := strings.Split("gomon .", " ")
	opts := NewOptions()
	dirArgs, cmdArgs := opts.Parse(args)
	if dirArgs[0] != "." {
		t.Fatal("dirArgs should be .")
	}
	if len(cmdArgs) != 0 {
		t.Fatal("cmdArgs should be empty")
	}
}

func TestParseCmdArg(t *testing.T) {
	args := strings.Split("gomon . -- go run -x server.go", " ")
	opts := NewOptions()
	dirArgs, cmdArgs := opts.Parse(args)
	if dirArgs[0] != "." {
		t.Fatal("dirArgs should be .")
	}
	if len(cmdArgs) != 4 ||
		strings.Join(cmdArgs, " ") != "go run -x server.go" {
		t.Fatal("Args after -- should be treated as cmdArgs")
	}
}
