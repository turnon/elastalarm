package notifiers

import (
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

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

	msg := gomail.NewMessage()
	msg.SetHeader("From", "1020715764@qq.com")
	msg.SetHeader("To", "1020715764@qq.com")
	msg.SetHeader("Subject", "liic测试")
	msg.SetBody("text/plain", "我是正文")

	d := gomail.NewPlainDialer("smtp.qq.com", 465, "1020715764@qq.com", "fkgjfk757jnuiry")
	if err := d.DialAndSend(msg); err != nil {
		panic(err)
	}
}
