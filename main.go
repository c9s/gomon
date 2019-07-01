package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/c9s/gomon/logger"
	"github.com/c9s/gomon/notify"
	"github.com/fsnotify/fsnotify"
)

var versionStr = "0.1.0"

var notifier notify.Notifier = nil

func main() {
	dirArgs, cmdArgs := options.Parse(os.Args)
	dirArgs = FilterExistPaths(dirArgs)

	var matchAll = false
	var alwaysNotify = false

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
		notify.InstallGrowlIcons()
		os.Exit(0)
		return
	}

	matchAll = options.Bool("matchall")
	alwaysNotify = options.Bool("alwaysnotify")

	// dynamically build the command list
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

	if options.Bool("d") {
		logger.Instance().SetLevel(logrus.DebugLevel)
	}

	if len(cmdArgs) > 0 {
		cmds.Add(Command(cmdArgs))
	} else if cmds.Len() == 0 {
		// default to go build
		cmds.Add(goCommands["build"])
	}

	if len(dirArgs) == 0 {
		var cwd, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dirArgs = []string{cwd}
	}

	if runtime.GOOS == "darwin" {
		logger.Infoln("Setting up Notification Center for OS X ...")
		notifier = notify.NewOSXNotifier()
	}
	if notifier == nil {
		if _, err := os.Stat("/Applications/Growl.app"); err == nil {
			logger.Infoln("Found Growl.app, setting up GNTP notifier...")
			notifier = notify.NewGNTPNotifier(options.String("gntp"), "gomon")
		}
	}
	if notifier == nil {
		notifier = notify.NewTextNotifier()
	}

	logger.Infoln("Watching", dirArgs, "for", cmds)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error(err)
		return
	}
	defer watcher.Close()

	for _, dir := range dirArgs {
		if options.Bool("R") {
			subfolders := Subfolders(dir)
			for _, f := range subfolders {
				err = watcher.Add(f)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			err = watcher.Add(dir)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	var jobRunner = &JobRunner{
		builder: &JobBuilder{
			// Job template arguments
			Commands:        cmds.commands,
			Args:            []string{},
			AppendFilename:  options.Bool("F"),
			ChangeDirectory: options.Bool("chdir"),
		},
	}

	var patternStr string = options.String("m")
	if len(patternStr) == 0 {
		// the empty regexp matches everything anyway
		matchAll = true
	}

	var pattern = regexp.MustCompile(patternStr)
	var once sync.Once
	var ctx = context.Background()

	for {
		select {

		case <-ctx.Done():
			return

		case err := <-watcher.Errors:
			log.Println("Error:", err)

		case e := <-watcher.Events:
			var matched = matchAll
			if !matched {
				matched = pattern.MatchString(e.Name)
			}

			if !matched {
				if options.Bool("d") {
					logger.Debugf("Ignored file=%s", e)
				}
				continue
			}

			if options.Bool("d") {
				logger.Debugf("Event=%+v", e)
			} else {
				switch e.Op {
				case fsnotify.Create:
					logger.Infoln("Create", e.Name)
				case fsnotify.Write:
					logger.Infoln("Write", e.Name)
				case fsnotify.Remove:
					logger.Infoln("Remove", e.Name)
				case fsnotify.Rename:
					logger.Infoln("Rename", e.Name)
				}
			}

			// TODO: time.ParseDuration
			// go fmt vim plugin will rename the file and then create a new file
			// In order to handle the batch operation, a delay is needed.
			go func(filename string) {
				once.Do(func() {
					var timer = time.After(500 * time.Millisecond)
					// duration to avoid to run commands frequency at once
					<-timer
					var duration, err = jobRunner.RunAndNotify(context.Background(), filename, alwaysNotify)
					if err != nil {
						logger.Errorf("Build failed: %v", err.Error())
					} else {
						logger.Infoln("Successful build:", duration)
					}
				})
				once = sync.Once{}
			}(e.Name)

		}
	}

}
