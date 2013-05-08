package main

import "github.com/howeyc/fsnotify"
import "log"
import "flag"
import "fmt"
import "os"
import "os/exec"
import "regexp"

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
	cmds []Command
}

func main() {
	var helpFlag = flag.Bool("h", false, "Show Help")
	var buildFlag = flag.Bool("b", false, "Run `go build`, the default behavior")
	var testFlag = flag.Bool("t", false, "Run `go test`")
	var installFlag = flag.Bool("i", false, "Run `go install`")
	var fmtFlag = flag.Bool("f", false, "Run `go fmt`")
	var runFlag = flag.Bool("r", false, "Run `go run`")
	var versionFlag = flag.Bool("v", false, "Version")
	var xFlag = flag.Bool("x", false, "Show verbose command")

	var useGrowl = flag.Bool("growl", false, "Use Growler")
	var gntpServer = flag.String("gntp", "", "The GNTP DSN")
	if *useGrowl && *gntpServer == "" {
		*gntpServer = "127.0.0.1:23053"
	}

	flag.Parse()
	args := flag.Args()

	if *helpFlag {
		fmt.Println("Usage: gomon [options] [dir] [-- command]")
		fmt.Println("   -b build")
		fmt.Println("   -t test")
		fmt.Println("   -i install")
		fmt.Println("   -x show verbose command")
		fmt.Println("   -h help")
		os.Exit(0)
	}
	if *versionFlag {
		fmt.Printf("gomon %s\n", versionStr)
		os.Exit(0)
	}

	var dirs = []string{}
	var cmds = []string{}

	if len(args) > 0 {
		var hasDash bool = false
		for _, a := range args {
			if a == "--" {
				hasDash = true
			}
		}
		if hasDash {
			var takeDir = true
			for _, a := range args {
				if takeDir && a == "--" {
					takeDir = false
					continue
				}
				if takeDir {
					dirs = append(dirs, a)
				} else {
					cmds = append(cmds, a)
				}
			}
		} else {
			// detect the argument (as path or command)
			var arePaths bool = true
			for _, a := range args {
				var exists, _ = FileExists(a)
				if !exists {
					arePaths = false
				}
			}
			if arePaths {
				dirs = args
			} else {
				cmds = args
			}
		}
	}

	if len(cmds) == 0 {
		if *testFlag {
			cmds = goCommands["test"]
		} else if *buildFlag {
			cmds = goCommands["build"]
		} else if *installFlag {
			cmds = goCommands["install"]
		} else if *fmtFlag {
			cmds = goCommands["fmt"]
		} else if *runFlag {
			cmds = goCommands["run"]
		} else {
			// default behavior
			cmds = goCommands["build"]
		}
		if *xFlag && len(cmds) > 0 {
			cmds = append(cmds, "-x")
		}
	}

	if len(cmds) == 0 {
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

	fmt.Println("Watching", dirs, "for", cmds)

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
	var cmd *exec.Cmd

	runCommand := func(cmd *exec.Cmd) {
		var err error
		err = cmd.Start()
		if err != nil {
			log.Println(err)
			if *useGrowl {
				notifyFail(gntpServer, err.Error(), "")
			}
			wasFailed = true
			return
		}
		err = cmd.Wait()
		if err != nil {
			log.Println(err)
			if *useGrowl {
				notifyFail(gntpServer, err.Error(), "")
			}
			wasFailed = true
			return
		}

		// fixed
		if wasFailed {
			wasFailed = false
			if *useGrowl {
				notifyFixed(gntpServer, "Fixed", "")
			}
			fmt.Println("Congratulations! It's fixed!")
		}
	}

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

			if cmd != nil && cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				err := cmd.Process.Kill()
				if err != nil {
					log.Println(err)
				}
			}
			cmd = exec.Command(cmds[0], cmds[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			go runCommand(cmd)

		case err := <-watcher.Error:
			log.Println("Error:", err)
		}
	}

	watcher.Close()
}
