package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var versionStr = "0.1.0"

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

		// everything after the dash, should be the command arguments
		if !hasDash && flag[0] == '-' {
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
			if !hasDash {
				if exists, _ := FileExists(arg); exists {
					dirArgs = append(dirArgs, arg)
				} else {
					log.Printf("Invalid path: '%v'", arg)
				}
			} else {
				cmdArgs = append(cmdArgs, arg)
			}
		}
	}

	if options.Bool("h") {
		fmt.Println("Usage: gomon [options] [dir] [-- command]")
		for _, option := range options {
			if _, ok := option.value.(string); ok {
				fmt.Printf("  -%s=%s: %s\n", option.flag, option.value, option.description)
			} else {
				fmt.Printf("  -%s: %s\n", option.flag, option.description)
			}
		}
		os.Exit(0)
	}
	if options.Bool("v") {
		fmt.Printf("gomon %s\n", versionStr)
		os.Exit(0)
	}

	if options.Bool("install-growl-icons") {
		installGrowlIcons()
		os.Exit(0)
	}

	var cmds = CommandSet{}
	var cmd = Command(cmdArgs)

	_ = cmds

	if len(cmd) == 0 {
		if options.Bool("t") {
			cmd = goCommands["test"]
		} else if options.Bool("i") {
			cmd = goCommands["install"]
		} else if options.Bool("f") {
			cmd = goCommands["fmt"]
		} else if options.Bool("r") {
			cmd = goCommands["run"]
		} else if options.Bool("b") {
			cmd = goCommands["build"]
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
			err = watcher.WatchFlags(f, fsnotify.FSN_ALL)
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
			failed("Failed!")
			wasFailed = true
			return
		}
		err = task.Wait()
		if err != nil {
			log.Println(err)
			if options.Bool("growl") {
				notifyFail(options.String("gntp"), err.Error(), "")
			}
			failed("Failed!")
			wasFailed = true
			return
		}

		// fixed
		if wasFailed {
			wasFailed = false
			if options.Bool("growl") {
				notifyFixed(options.String("gntp"), "Congratulations! It's fixed!", "")
			}
			success("Congratulations! It's fixed!")
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
				if options.Bool("d") {
					log.Println("Ignore:", e)
				}
				continue
			}

			if options.Bool("d") {
				log.Println("Event:", e)
			} else {
				log.Println(e.Name)
			}

			if !fired {
				fired = true
				go func() {
					// duration to avoid to run commands frequency at once
					select {
					case <-time.After(100 * time.Millisecond):
						fired = false
						if task != nil && task.ProcessState != nil && !task.ProcessState.Exited() {
							fmt.Println("Stopping Task...")
							err := task.Process.Kill()
							if err != nil {
								log.Println(err)
							}
						}
						fmt.Println("Running Task:", cmd)
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
