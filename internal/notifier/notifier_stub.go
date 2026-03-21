//go:build !darwin

package notifier

type noop struct{}

func New() Sender {
	return noop{}
}

func (noop) Notify(title, message string) error {
	return nil
}
