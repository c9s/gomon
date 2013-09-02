package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/mattn/go-gntp"
)

// var server = flag.String("s", "127.0.0.1:23053", "GNTP server")
// var action = flag.String("a", "", "Click action")
//	var buf bytes.Buffer
//	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
//	cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
//	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)
//	err := cmd.Run()

func createNotification(server string) *gntp.Client {
	growl := gntp.NewClient()
	// defualt GNTP Server
	growl.Server = server
	growl.AppName = "gomon"
	growl.Register([]gntp.Notification{
		gntp.Notification{
			Event:   "success",
			Enabled: false,
		}, gntp.Notification{
			Event:   "failed",
			Enabled: true,
		},
	})
	return growl
}

func notifyFixed(server string, text, callback string) {
	growl := createNotification(server)
	growl.Notify(&gntp.Message{
		Event:    "success",
		Title:    "Fixed",
		Text:     text,
		Callback: callback,
		Icon:     icon("success"),
	})
}

func notifyFail(server string, text, callback string) {
	growl := createNotification(server)
	growl.Notify(&gntp.Message{
		Event:    "failed",
		Title:    "Failed",
		Text:     text,
		Callback: callback,
		Icon:     icon("failed"),
	})
}

func success(msg string) {
	ct.ChangeColor(ct.Black, true, ct.Green, true)
	fmt.Print(msg)
	ct.ResetColor()
	fmt.Println()
}

func failed(msg string) {
	ct.ChangeColor(ct.Black, true, ct.Red, true)
	fmt.Print(msg)
	ct.ResetColor()
	fmt.Println()
}
