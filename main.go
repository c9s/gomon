package main

import "github.com/howeyc/fsnotify"
import "log"
import "fmt"
import "os"
import "os/exec"
import "regexp"
import "strings"
import "time"

var versionStr = "0.1.0"

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

type gomonOption struct {
	flag        string
	value       interface{}
	description string
}

type gomonOptions []*gomonOption

var options = gomonOptions{
	{"h", false, "Show Help"},
	{"b", true, "Run `go build`, the default behavior"},
	{"t", false, "Run `go test`"},
	{"i", false, "Run `go install`"},
	{"f", false, "Run `go fmt`"},
	{"r", false, "Run `go run`"},
	{"x", false, "Show verbose command"},
	{"v", false, "Show version"},
	{"growl", false, "Use Growler"},
	{"gntp", "127.0.0.1:23053", "The GNTP DSN"},
}

func (options gomonOptions) Get(flag string) *gomonOption {
	for _, option := range options {
		if option.flag == flag {
			return option
		}
	}
	return nil
}

func (options gomonOptions) String(flag string) string {
	for _, option := range options {
		if option.flag == flag {
			s, _ := option.value.(string)
			return s
		}
	}
	return ""
}

func (options gomonOptions) Bool(flag string) bool {
	for _, option := range options {
		if option.flag == flag {
			b, _ := option.value.(bool)
			return b
		}
	}
	return false
}

func main() {
	var dirArgs = []string{}
	var cmdArgs = []string{}

	var hasDash bool = false
	for n := 1; n < len(os.Args); n++ {
		arg := os.Args[n]
		if arg == "--" {
			hasDash = true
			continue
		}
		tokens := strings.SplitN(arg, "=", 2)
		flag, value := "", ""
		switch len(tokens) {
		case 1:
			flag = tokens[0]
		case 2:
			flag = tokens[0]
			value = tokens[1]
		default:
			continue
		}
		if flag[0] == '-' {
			option := options.Get(flag[1:])
			if option == nil {
				log.Printf("Invalid option: '%v'\n", flag)
			} else {
				if _, ok := option.value.(string); ok {
					option.value = value
				} else {
					option.value = true
				}
			}
		} else {
			if hasDash {
				if exists, _ := FileExists(arg); exists {
					dirArgs = append(dirArgs, arg)
				} else {
					log.Printf("Invalid path are specified: '%v'", arg)
				}
			} else {
				cmdArgs = append(cmdArgs, arg)
			}
		}
	}

	if options.Bool("h") {
		fmt.Println("Usage: gomon [options] [dir] [-- command]")
		for _, option := range options {
			fmt.Printf("  -%s: %s\n", option.flag, option.description)
		}
		os.Exit(0)
	}
	if options.Bool("v") {
		fmt.Printf("gomon %s\n", versionStr)
		os.Exit(0)
	}

	var cmds = CommandSet{}
	var cmd = Command(cmdArgs)

	_ = cmds

	if len(cmd) == 0 {
		if options.Bool("t") {
			cmd = goCommands["test"]
		} else if options.Bool("b") {
			cmd = goCommands["build"]
		} else if options.Bool("i") {
			cmd = goCommands["install"]
		} else if options.Bool("f") {
			cmd = goCommands["fmt"]
		} else if options.Bool("r") {
			cmd = goCommands["run"]
		} else {
			// default behavior
			cmd = goCommands["build"]
		}
		if options.Bool("x") && len(cmd) > 0 {
			cmd = append(cmd, "-x")
		}
	}

	if len(cmd) == 0 {
		fmt.Println("No command specified")
		os.Exit(2)
	}

	if len(dirArgs) == 0 {
		var cwd, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dirArgs = []string{cwd}
	}

	fmt.Println("Watching", dirArgs, "for", cmd)

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range dirArgs {
		subfolders := Subfolders(dir)
		for _, f := range subfolders {
			err = watcher.WatchFlags(f, fsnotify.FSN_CREATE|fsnotify.FSN_MODIFY)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	var wasFailed bool = false
	var task *exec.Cmd

	runCommand := func(task *exec.Cmd) {
		var err error
		err = task.Start()
		if err != nil {
			log.Println(err)
			if options.Bool("growl") {
				notifyFail(options.String("gntp"), err.Error(), "")
			}
			wasFailed = true
			return
		}
		err = task.Wait()
		if err != nil {
			log.Println(err)
			if options.Bool("growl") {
				notifyFail(options.String("gntp"), err.Error(), "")
			}
			wasFailed = true
			return
		}

		// fixed
		if wasFailed {
			wasFailed = false
			if options.Bool("growl") {
				notifyFixed(options.String("gntp"), "Congratulations!", "")
			}
			fmt.Println("Congratulations! It's fixed!")
		}
	}

	var fired bool = false
	for {
		select {
		case e := <-watcher.Event:
			matched, err := regexp.MatchString("\\.(go|c|h)$", e.Name)
			if err != nil {
				log.Println(err)
			}

			if !matched {
				continue
			}

			log.Println("Event:", e)

			if !fired {
				fired = true
				go func() {
					// duration to avoid to run commands frequency at once
					select {
					case <-time.After(100 * time.Millisecond):
						fired = false
						if task != nil && task.ProcessState != nil && !task.ProcessState.Exited() {
							err := task.Process.Kill()
							if err != nil {
								log.Println(err)
							}
						}
						task = exec.Command(cmd[0], cmd[1:]...)
						task.Stdout = os.Stdout
						task.Stderr = os.Stderr
						runCommand(task)
					}
				}()
			}

		case err := <-watcher.Error:
			log.Println("Error:", err)
		}
	}

	watcher.Close()
}
