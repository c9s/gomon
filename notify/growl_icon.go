package notify

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var iconDirectory string

func getConfigDir() string {
	u, err := user.Current()
	if err != nil {
		log.Fatal("failed to get home directory: ", err)
	}
	return filepath.Join(u.HomeDir, ".gomon")
}

func getIconDir() string {
	if iconDirectory != "" {
		return iconDirectory
	}
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
	} else {
		dir = getConfigDir()
	}
	iconDirectory = filepath.Join(dir, "gomon")
	return iconDirectory
}

func fileExists(path string) (bool, error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return false, err
	}
	_, err = file.Stat()
	if err != nil {
		return false, err
	}
	err = file.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

func icon(name string) string {
	f := filepath.Join(getIconDir(), name+".png")
	if _, err := fileExists(f); err == nil {
		return f
	}
	return ""
}

func download(target, path string) {
	r, err := http.Get(target)
	if err != nil {
		log.Fatal("failed to download file: ", err)
	}
	defer r.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		log.Fatal("failed to create file: ", err)
	}
	defer out.Close()
	io.Copy(out, r.Body)
	fmt.Printf("Downloaded %s into %s\n", target, path)
}

// InstallGrowlIcons install icon images for Growl Notification
func InstallGrowlIcons() {
	dir := getIconDir()
	_, err := os.Stat(dir)
	if err != nil {
		if os.MkdirAll(dir, 0700) != nil {
			log.Fatal("failed to create directory: ", err)
		}
	}
	download(
		"https://raw.github.com/c9s/gomon/gh-pages/icons/success.png",
		filepath.Join(getIconDir(), "success.png"))
	download(
		"https://raw.github.com/c9s/gomon/gh-pages/icons/failed.png",
		filepath.Join(getIconDir(), "failed.png"))
}
