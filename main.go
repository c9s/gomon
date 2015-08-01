package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/howeyc/fsnotify"
)

var versionStr = "0.1.0"

func main() {
	dirArgs, cmdArgs := options.Parse(os.Args)
	dirArgs = FilterExistPaths(dirArgs)

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

	matchAll := false
	if options.Bool("matchall") {
		matchAll = true
	}

	var cmds = CommandList{}
	if options.Bool("f") {
		cmds.Add(goCommands["fmt"])
	}
	if options.Bool("t") {
		cmds.Add(goCommands["test"])
	}
	if options.Bool("b") {
		cmds.Add(goCommands["build"])
	}
	if options.Bool("r") {
		cmds.Add(goCommands["run"])
	}
	if options.Bool("i") {
		cmds.Add(goCommands["install"])
	}
	if options.Bool("x") {
		cmds.AppendOption("-x")
	}
	if len(cmdArgs) > 0 {
		cmds.Add(Command(cmdArgs))
	} else {
		if cmds.Len() == 0 {
			cmds.Add(goCommands["build"])
		}
	}

	if len(dirArgs) == 0 {
		var cwd, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dirArgs = []string{cwd}
	}

	fmt.Println("Watching", dirArgs, "for", cmds)

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range dirArgs {
		if options.Bool("R") {
			subfolders := Subfolders(dir)
			for _, f := range subfolders {
				err = watcher.WatchFlags(f, fsnotify.FSN_ALL)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			err = watcher.WatchFlags(dir, fsnotify.FSN_ALL)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	var wasFailed bool = false

	runCommand := func(filename string) {
		var dirOpt *string
		var dir = filepath.Dir(filename)
		if options.Bool("chdir") {
			dirOpt = &dir
		} else {
			dirOpt = nil
		}
		if options.Bool("F") {
			cmds.SetFilename(filename)
		} else {
			cmds.ClearFilename()
		}
		err := cmds.Run(dirOpt)
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

	var pattern string = options.String("m")
	if len(pattern) == 0 {
		// the empty regexp matches everything anyway
		matchAll = true
	}
	var fired bool = false

	for {
		select {
		case e := <-watcher.Event:
			var err error
			matched := matchAll
			if !matched {
				matched, err = regexp.MatchString(pattern, e.Name)
				if err != nil {
					log.Println(err)
				}
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
				go func(filename string) {
					// duration to avoid to run commands frequency at once
					select {
					case <-time.After(200 * time.Millisecond):
						err := cmds.StopTask()
						if err != nil {
							log.Println(err)
						}
						fmt.Println("Running Task:", cmds)
						runCommand(filename)
						fired = false
					}
				}(e.Name)
			}

		case err := <-watcher.Error:
			log.Println("Error:", err)
		}
	}

	watcher.Close()
}
