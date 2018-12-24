package notifiers

import "fmt"

// Email to send
type Email struct {
	title string
	body  string
}

// SetTitle set email title
func (m *Email) SetTitle(s string) {
	m.title = s
}

// SetBody set email body
func (m *Email) SetBody(s string) {
	m.body = s
}

// Notify send email
func (m *Email) Notify() {
	fmt.Println(m.title, m.body)
}
