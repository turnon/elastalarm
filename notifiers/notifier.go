package notifiers

// Notifier interface
type Notifier interface {
	SetTitle(s string)
	SetBody(s string)
	Notify()
}
