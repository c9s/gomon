package main

import "os"

func FileExists(path string) (bool, error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return false, err
	}
	_, err = file.Stat()
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
	return stat.IsDir(), nil
}
