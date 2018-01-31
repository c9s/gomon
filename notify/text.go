package notify

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
)

type Notifier interface {
	NotifySucceeded(text string) error
	NotifyFailed(text string) error
	NotifyFixed(text string) error
}

type TextNotifier struct {
}

func NewTextNotifier() *TextNotifier {
	return &TextNotifier{}
}

func (n *TextNotifier) NotifySucceeded(msg string) error {
	ct.ChangeColor(ct.Black, false, ct.Green, true)
	fmt.Print(msg)
	ct.ResetColor()
	_, err := fmt.Println()
	return err
}

func (n *TextNotifier) NotifyFixed(msg string) error {
	ct.ChangeColor(ct.Black, false, ct.Green, true)
	fmt.Print(msg)
	ct.ResetColor()
	_, err := fmt.Println()
	return err
}

func (n *TextNotifier) NotifyFailed(msg string) error {
	ct.ChangeColor(ct.Black, false, ct.Red, true)
	fmt.Print(msg)
	ct.ResetColor()
	_, err := fmt.Println()
	return err
}
