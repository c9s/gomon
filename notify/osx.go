// +build darwin

package notify

import (
	gosxnotifier "github.com/deckarep/gosx-notifier"
)

type OSXNotifier struct{}

func NewOSXNotifier() *OSXNotifier {
	return &OSXNotifier{}
}

// NotifySucceeded show notification of succeeded
func (n *OSXNotifier) NotifySucceeded(title, subtitle string) error {
	note := gosxnotifier.NewNotification("Build Succeeded") // message body
	note.Title = title
	note.Subtitle = subtitle
	note.Sound = gosxnotifier.Sound("Submarine")
	return note.Push()
}

// NotifyFixed show notification of fixed
func (n *OSXNotifier) NotifyFailed(title, subtitle string) error {
	note := gosxnotifier.NewNotification("Build Failed") // message body
	note.Title = title
	note.Subtitle = subtitle
	note.Sound = gosxnotifier.Basso
	return note.Push()
}

// NotifyFailed show notification of failed
func (n *OSXNotifier) NotifyFixed(title, subtitle string) error {
	note := gosxnotifier.NewNotification("Build Fixed") // message body
	note.Title = title
	note.Subtitle = subtitle
	note.Sound = gosxnotifier.Sound("Submarine")
	return note.Push()
}
