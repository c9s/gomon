// +build !darwin

package notify

func NewOSXNotifier() Notifier {
	return nil
}
