package main

import (
	"github.com/mattn/go-gntp"
)

// var server = flag.String("s", "127.0.0.1:23053", "GNTP server")
// var action = flag.String("a", "", "Click action")
//	var buf bytes.Buffer
//	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
//	cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
//	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)
//	err := cmd.Run()

type GNTPNotifier struct {
	*gntp.Client
}

func NewGNTPNotifier(server string, appName string) *GNTPNotifier {
	client := gntp.NewClient()
	// defualt GNTP Server
	client.Server = server
	client.AppName = appName
	client.Register([]gntp.Notification{
		gntp.Notification{
			Event:   "succeeded",
			Enabled: true,
		}, gntp.Notification{
			Event:   "failed",
			Enabled: true,
		},
	})
	return &GNTPNotifier{client}
}

func (n *GNTPNotifier) NotifySucceeded(text string) error {
	return n.Client.Notify(&gntp.Message{
		Event: "succeeded",
		Title: "Succeeded",
		Text:  text,
		Icon:  icon("success"),
	})
}

func (n *GNTPNotifier) NotifyFixed(text string) error {
	return n.Client.Notify(&gntp.Message{
		Event: "succeeded",
		Title: "Fixed",
		Text:  text,
		Icon:  icon("success"),
	})
}

func (n *GNTPNotifier) NotifyFailed(text string) error {
	return n.Client.Notify(&gntp.Message{
		Event: "failed",
		Title: "Failed",
		Text:  text,
		Icon:  icon("failed"),
	})
}
