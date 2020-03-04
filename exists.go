package main

import (
	"log"
	"os"
	"path/filepath"
)

//
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

func isDir(path string) (bool, error) {
	return dirExists(path)
}

func dirExists(path string) (bool, error) {
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

func subfolders(path string) (paths []string) {
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
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

func filterExistPaths(paths []string) []string {
	var result []string
	for _, path := range paths {
		if exists, _ := fileExists(path); exists {
			result = append(result, path)
		} else {
			log.Printf("Invalid path: '%v'", path)
		}
	}
	return result
}
