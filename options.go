package main

import (
	"log"
	"strings"
)

type Option struct {
	flag        string
	value       interface{}
	description string
}

type Options []*Option

func NewOptions() Options {
	return Options{
		{"h", false, "Show Help"},
		{"b", false, "Run `go build`, the default behavior"},
		{"t", false, "Run `go test`"},
		{"i", false, "Run `go install`"},
		{"f", false, "Run `go fmt`"},
		{"m", `\.(go|c|h)$`, "Pattern to match"},
		{"r", false, "Run `go run`"},
		{"x", false, "Show verbose command"},
		{"v", false, "Show version"},
		{"F", false, "Append modified filename to command"},
		{"d", false, "Print debug message"},
		{"R", false, "Watch directory recursively"},
		{"matchall", false, "Match all files (equivalent to -m='')"},
		{"alwaysnotify", false, "Always send notification"},
		{"notify", "", "Select notifier: osx / growl / off"},
		{"chdir", false, "Run commands on directory"},
		{"growl", false, "Use Growler"},
		{"install-growl-icons", false, "Install growl icons"},
		{"gntp", "127.0.0.1:23053", "The GNTP DSN"},
	}
}

var options = NewOptions()

func (options Options) Has(flag string) bool {
	for _, option := range options {
		if option.flag == flag {
			return true
		}
	}
	return false
}

func (options Options) Get(flag string) *Option {
	for _, option := range options {
		if option.flag == flag {
			return option
		}
	}
	return nil
}

func (options Options) String(flag string) string {
	for _, option := range options {
		if option.flag == flag {
			s, _ := option.value.(string)
			return s
		}
	}
	return ""
}

func (options Options) Bool(flag string) bool {
	for _, option := range options {
		if option.flag == flag {
			b, _ := option.value.(bool)
			return b
		}
	}
	return false
}

func (options Options) IsBool(flag string) bool {
	for _, option := range options {
		if option.flag == flag {
			_, ok := option.value.(bool)
			return ok
		}
	}
	return false
}

func (options Options) Parse(args []string) (dirArgs []string, cmdArgs []string) {
	dirArgs = make([]string, 0)
	cmdArgs = make([]string, 0)

	var hasDash bool = false
	var nArgs = len(args)
	for n := 1; n < nArgs; n++ {
		arg := args[n]
		if arg == "--" {
			hasDash = true
			continue
		}
		if hasDash {
			// everything after the dash, should be the command arguments
			cmdArgs = append(cmdArgs, arg)
		} else {
			if arg[0] == '-' {
				tokens := strings.SplitN(arg, "=", 2)
				flag, value := "", ""
				switch len(tokens) {
				case 1:
					flag = tokens[0]
					if n < nArgs-1 && !options.Has(flag) && !options.IsBool(flag[1:]) {
						value = args[n+1]
						n++
					}
				case 2:
					flag = tokens[0]
					value = tokens[1]
				default:
					continue
				}

				option := options.Get(flag[1:])
				if option == nil {
					log.Fatalf("Invalid option: '%v'\n", flag)
				} else {
					if _, ok := option.value.(string); ok {
						option.value = value
					} else {
						option.value = true
					}
				}
			} else {
				dirArgs = append(dirArgs, arg)
			}
		}
	}
	return
}
