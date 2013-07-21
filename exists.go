package main

import (
	"log"
	"os"
	"path/filepath"
)

func FileExists(path string) (bool, error) {
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

func IsDir(path string) (bool, error) {
	return DirExists(path)
}

func DirExists(path string) (bool, error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return false, err
	}
	stat, err := file.Stat()
	if err != nil {
		return false, err
	}
	err = file.Close()
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}

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

func FilterExistPaths(paths []string) []string {
	var result []string
	for _, path := range paths {
		if exists, _ := FileExists(path); exists {
			result = append(result, path)
		} else {
			log.Printf("Invalid path: '%v'", path)
		}
	}
	return result
}
