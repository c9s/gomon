package main

import "github.com/howeyc/fsnotify"
import "log"
import "flag"
import "fmt"
import "os"
import "os/exec"
import "regexp"
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

func main() {
	var helpFlag = flag.Bool("h", false, "Show Help")
	var buildFlag = flag.Bool("b", false, "Run `go build`, the default behavior")
	var testFlag = flag.Bool("t", false, "Run `go test`")
	var installFlag = flag.Bool("i", false, "Run `go install`")
	var fmtFlag = flag.Bool("f", false, "Run `go fmt`")
	var runFlag = flag.Bool("r", false, "Run `go run`")
	// var allFlag = flag.Bool("a", false, "Run build, test, fmt and install")

	var versionFlag = flag.Bool("v", false, "Version")
	var xFlag = flag.Bool("x", false, "Show verbose command")

	var useGrowl = flag.Bool("growl", false, "Use Growler")
	var gntpServer = flag.String("gntp", "127.0.0.1:23053", "The GNTP DSN")

	flag.Parse()

	var args = flag.Args()
	if *helpFlag {
		fmt.Println("Usage: gomon [options] [dir] [-- command]")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *versionFlag {
		fmt.Printf("gomon %s\n", versionStr)
		os.Exit(0)
	}

	var dirs = []string{}
	var cmds = CommandSet{}
	var cmd = Command{}

	_ = cmds

	var hasDash bool = false
	for _, a := range args {
		if a == "--" {
			hasDash = true
			continue
		}
		if !hasDash {
			if exists, _ := FileExists(a); exists {
				dirs = append(dirs, a)
			} else {
				log.Printf("Invalid path are specified: '%v'", a)
			}
		} else {
			cmd = append(cmd, a)
		}
	}

	if len(cmd) == 0 {
		if *testFlag {
			cmd = goCommands["test"]
		} else if *buildFlag {
			cmd = goCommands["build"]
		} else if *installFlag {
			cmd = goCommands["install"]
		} else if *fmtFlag {
			cmd = goCommands["fmt"]
		} else if *runFlag {
			cmd = goCommands["run"]
		} else {
			// default behavior
			cmd = goCommands["build"]
		}
		if *xFlag && len(cmd) > 0 {
			cmd = append(cmd, "-x")
		}
	}

	if len(cmd) == 0 {
		fmt.Println("No command specified")
		os.Exit(2)
	}

	if len(dirs) == 0 {
		var cwd, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dirs = []string{cwd}
	}

	fmt.Println("Watching", dirs, "for", cmd)

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range dirs {
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
			if *useGrowl {
				notifyFail(*gntpServer, err.Error(), "")
			}
			wasFailed = true
			return
		}
		err = task.Wait()
		if err != nil {
			log.Println(err)
			if *useGrowl {
				notifyFail(*gntpServer, err.Error(), "")
			}
			wasFailed = true
			return
		}

		// fixed
		if wasFailed {
			wasFailed = false
			if *useGrowl {
				notifyFixed(*gntpServer, "Fixed", "")
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
