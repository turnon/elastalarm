package notifiers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	gomail "gopkg.in/gomail.v2"
)

// Notifier interface
type Notifier interface {
	SetTitle(s string)
	SetBody(s string)
	Notify()
}

type Msg struct {
	Title, Body *string
}

var Names = make(map[string](func(*Msg)))

func init() {
	Names["stdout"] = stdout
	Names["email"] = emailFunc()
}

func stdout(m *Msg) {
	fmt.Println(*m.Title, *m.Body)
}

func emailFunc() func(*Msg) {
	server := os.Getenv("ESALARM_MAIL_SERVER")
	if server == "" {
		return nil
	}

	hostPort := strings.Split(server, ":")
	host := hostPort[0]
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		panic(err)
	}

	from := os.Getenv("ESALARM_MAIL_FROM")
	passwd := os.Getenv("ESALARM_MAIL_PASSWD")
	to := os.Getenv("ESALARM_MAIL_TO")

	return func(m *Msg) {
		msg := gomail.NewMessage()
		msg.SetHeader("From", from)
		msg.SetHeader("To", to)
		msg.SetHeader("Subject", *m.Title)
		msg.SetBody("text/plain", *m.Body)

		d := gomail.NewPlainDialer(host, port, from, passwd)
		if err := d.DialAndSend(msg); err != nil {
			panic(err)
		}
	}
}
