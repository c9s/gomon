// +build !windows

package main

import (
	"fmt"
	"github.com/koyachi/go-term-ansicolor/ansicolor"
)

func success(msg string) {
	fmt.Println(ansicolor.Black(ansicolor.OnGreen(msg)))
}

func failed(msg string) {
	fmt.Println(ansicolor.Black(ansicolor.OnRed(msg)))
}
