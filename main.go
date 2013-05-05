package main
import "github.com/howeyc/fsnotify"
import "log"
import "flag"
import "fmt"
import "path/filepath"
import "os"
import "os/exec"
import "regexp"







func Subfolders(path string) (paths []string) {
    filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            name := info.Name()
            // skip folders that begin with a dot
            hidden := filepath.HasPrefix(name, ".") && name != "." && name != ".."
            if hidden {
                return filepath.SkipDir
            } else {
                paths = append(paths, newPath)
            }
        }
        return nil
    })
    return paths
}

func main() {
	flag.Parse()
	args := flag.Args()

	var dirs = []string{}
	var cmds = []string{}

	var takeDir = true
	for _, a := range args {
		if a == "--" {
			takeDir = false
			continue
		}
		if takeDir {
			dirs = append(dirs, a)
		} else {
			cmds = append(cmds, a)
		}
	}

	fmt.Println( "Watching" ,  dirs , "for" , cmds )

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

	var cmd *exec.Cmd

	for {
		select {
		case e := <-watcher.Event:
			matched, err := regexp.MatchString("\\.go$", e.Name)
			if err != nil {
				log.Println(err)
			}

			if ! matched {
				continue
			}

			log.Println("Event:", e)

			if cmd != nil {
				err := cmd.Process.Kill()
				if err != nil {
					log.Println(err)
				}
			}
			cmd = exec.Command(cmds[0], cmds[1:]... )
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		case err := <-watcher.Error:
			log.Println("Error:", err)
		}
	}

    watcher.Close()
}
