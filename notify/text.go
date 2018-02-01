package notify

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
)

type Notifier interface {
	NotifySucceeded(text, subtitle string) error
	NotifyFailed(text, subtitle string) error
	NotifyFixed(text, subtitle string) error
}

type TextNotifier struct {
}

func NewTextNotifier() *TextNotifier {
	return &TextNotifier{}
}

func (n *TextNotifier) NotifySucceeded(title string, subtitle string) error {
	ct.ChangeColor(ct.Black, false, ct.Green, true)
	fmt.Print(title, ":", subtitle)

	ct.ResetColor()
	_, err := fmt.Println()
	return err
}

func (n *TextNotifier) NotifyFixed(title string, subtitle string) error {
	ct.ChangeColor(ct.Black, false, ct.Green, true)
	fmt.Print(title, ":", subtitle)

	ct.ResetColor()
	_, err := fmt.Println()
	return err
}

func (n *TextNotifier) NotifyFailed(title string, subtitle string) error {
	ct.ChangeColor(ct.Black, false, ct.Red, true)
	fmt.Print(title, ":", subtitle)
	ct.ResetColor()
	_, err := fmt.Println()
	return err
}
