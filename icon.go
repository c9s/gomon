package main

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

func iconDir() string {
	if iconDirectory != "" {
		return iconDirectory
	}
	u, err := user.Current()
	if err != nil {
		log.Fatal("failed to get home directory: ", err)
	}
	dir := filepath.Join(u.HomeDir, ".config")
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
	}
	iconDirectory = filepath.Join(dir, "gomon")
	return iconDirectory
}

func icon(name string) string {
	f := filepath.Join(iconDir(), name+".png")
	if _, err := FileExists(f); err == nil {
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

func installGrowlIcons() {
	dir := iconDir()
	_, err := os.Stat(dir)
	if err != nil {
		if os.Mkdir(dir, 0700) != nil {
			log.Fatal("failed to create directory: ", err)
		}
	}
	download(
		"https://raw.github.com/c9s/gomon/gh-pages/icons/success.png",
		filepath.Join(iconDir(), "success.png"))
	download(
		"https://raw.github.com/c9s/gomon/gh-pages/icons/failed.png",
		filepath.Join(iconDir(), "failed.png"))
}

