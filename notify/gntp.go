package notify

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

// GNTPNotifier is
type GNTPNotifier struct {
	*gntp.Client
}

// NewGNTPNotifier return new instance of GNTPNotifier
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

// NotifySucceeded show notification of succeeded
func (n *GNTPNotifier) NotifySucceeded(title, subtitle string) error {
	return n.Client.Notify(&gntp.Message{
		Event: "succeeded",
		Title: title,
		Text:  subtitle,
		Icon:  icon("success"),
	})
}

// NotifyFixed show notification of fixed
func (n *GNTPNotifier) NotifyFixed(title, subtitle string) error {
	return n.Client.Notify(&gntp.Message{
		Event: "succeeded",
		Title: title,
		Text:  subtitle,
		Icon:  icon("success"),
	})
}

// NotifyFailed show notification of failed
func (n *GNTPNotifier) NotifyFailed(title, subtitle string) error {
	return n.Client.Notify(&gntp.Message{
		Event: "failed",
		Title: title,
		Text:  subtitle,
		Icon:  icon("failed"),
	})
}
