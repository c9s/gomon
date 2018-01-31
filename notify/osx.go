// +build darwin,!cgo
package notify

import (
	"github.com/deckarep/gosx-notifier"
)

type OSXNotifier struct{}

func NewOSXNotifier() *OSXNotifier {
	return &OSXNotifier{}
}

func (n *OSXNotifier) NotifySucceeded(title string) error {
	note := gosxnotifier.NewNotification(title)
	note.Title = title
	note.Subtitle = ""
	note.Sound = gosxnotifier.Sound("Submarine")
	return note.Push()
}

func (n *OSXNotifier) NotifyFailed(title string) error {
	note := gosxnotifier.NewNotification(title)
	note.Title = title
	note.Subtitle = ""
	note.Sound = gosxnotifier.Basso
	return note.Push()
}

func (n *OSXNotifier) NotifyFixed(title string) error {
	note := gosxnotifier.NewNotification(title)
	note.Title = title
	note.Subtitle = ""
	note.Sound = gosxnotifier.Sound("Submarine")
	return note.Push()
}
