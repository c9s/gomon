package main

import (
	"fmt"
	"github.com/anschelsc/doscolor"
	"os"
)

var cout = doscolor.NewWrapper(os.Stdout)

func success(msg string) {
	cout.Save()
	cout.Set(doscolor.Green | doscolor.Bright)
	fmt.Println(msg)
	cout.Restore()
}

func failed(msg string) {
	cout.Save()
	cout.Set(doscolor.Red | doscolor.Bright)
	fmt.Println(msg)
	cout.Restore()
}
