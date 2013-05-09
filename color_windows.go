package main

import (
	"fmt"
	"github.com/anschelsc/doscolor"
	"os"
)

var cout = doscolor.NewWrapper(os.Stdout)

func success(msg string) {
	cout.Set(doscolor.Green | doscolor.Bright)
	fmt.Println(msg)
	cout.Restore()
}
