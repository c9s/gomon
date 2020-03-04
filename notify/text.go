package notify

import (
	"fmt"

	ct "github.com/daviddengcn/go-colortext"
)

// Notifier is
type Notifier interface {
	NotifySucceeded(text, subtitle string) error
	NotifyFailed(text, subtitle string) error
	NotifyFixed(text, subtitle string) error
}

// TextNotifier is
type TextNotifier struct {
}

// NewTextNotifier return new instance of TextNotifier
func NewTextNotifier() *TextNotifier {
	return &TextNotifier{}
}

// NotifySucceeded show notification of succeeded
func (n *TextNotifier) NotifySucceeded(title string, subtitle string) error {
	ct.ChangeColor(ct.Black, false, ct.Green, true)
	fmt.Print(title, ":", subtitle)

	ct.ResetColor()
	_, err := fmt.Println()
	return err
}

// NotifyFixed show notification of fixed
func (n *TextNotifier) NotifyFixed(title string, subtitle string) error {
	ct.ChangeColor(ct.Black, false, ct.Green, true)
	fmt.Print(title, ":", subtitle)

	ct.ResetColor()
	_, err := fmt.Println()
	return err
}

// NotifyFailed show notification of failed
func (n *TextNotifier) NotifyFailed(title string, subtitle string) error {
	ct.ChangeColor(ct.Black, false, ct.Red, true)
	fmt.Print(title, ":", subtitle)
	ct.ResetColor()
	_, err := fmt.Println()
	return err
}
