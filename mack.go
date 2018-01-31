// +build darwin,!cgo
package main

import "github.com/everdev/mack"

func foo() {
	mack.Say("Starting process")
	// do stuff
	mack.Notify("Complete")
}
