package main

import (
	// "bytes"
	"github.com/mattn/go-gntp"
)

// var server = flag.String("s", "127.0.0.1:23053", "GNTP server")
// var action = flag.String("a", "", "Click action")
//	var buf bytes.Buffer
//	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
//	cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
//	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)
//	err := cmd.Run()

func createNotification(server *string) *gntp.Client {
	growl := gntp.NewClient()
	// defualt GNTP Server
	if server != nil {
		growl.Server = *server
	} else {
		growl.Server = "127.0.0.1:23053"
	}
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

func notifyFail(text, callback string, server *string) {
	growl := createNotification(server)
	growl.Notify(&gntp.Message{
		Event:    "failed",
		Title:    "Failed",
		Text:     text,
		Callback: callback,
	})
}
