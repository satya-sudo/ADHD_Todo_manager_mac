package notifier

type Sender interface {
	Notify(title, message string) error
}
