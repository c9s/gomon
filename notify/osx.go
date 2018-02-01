// +build darwin,!cgo
package notify

import (
	"github.com/deckarep/gosx-notifier"
)

type OSXNotifier struct{}

func NewOSXNotifier() *OSXNotifier {
	return &OSXNotifier{}
}

func (n *OSXNotifier) NotifySucceeded(title, subtitle string) error {
	note := gosxnotifier.NewNotification("Build Succeeded") // message body
	note.Title = title
	note.Subtitle = subtitle
	note.Sound = gosxnotifier.Sound("Submarine")
	return note.Push()
}

func (n *OSXNotifier) NotifyFailed(title, subtitle string) error {
	note := gosxnotifier.NewNotification("Build Failed") // message body
	note.Title = title
	note.Subtitle = subtitle
	note.Sound = gosxnotifier.Basso
	return note.Push()
}

func (n *OSXNotifier) NotifyFixed(title, subtitle string) error {
	note := gosxnotifier.NewNotification("Build Fixed") // message body
	note.Title = title
	note.Subtitle = subtitle
	note.Sound = gosxnotifier.Sound("Submarine")
	return note.Push()
}
