package notifiers

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hugozhu/godingtalk"
	"github.com/pkg/errors"
	gomail "gopkg.in/gomail.v2"
)

var Names = make(map[string](func(*Msg) error))

func init() {
	Names["stdout"] = stdout
	Names["email"] = emailFunc()
	Names["ding"] = dingFunc()
}

type Msg struct {
	Title, Body *string
}

func (msg *Msg) join(seperate string) string {
	return *msg.Title + seperate + *msg.Body
}

func stdout(m *Msg) error {
	fmt.Println(*m.Title, "\n\n", *m.Body)
	return nil
}

func emailFunc() func(*Msg) error {
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
	skipVerify := os.Getenv("ESALARM_MAIL_SKIP_VERIFY") != ""

	return func(m *Msg) error {
		msg := gomail.NewMessage()
		msg.SetHeader("From", from)
		msg.SetHeader("To", to)
		msg.SetHeader("Subject", *m.Title)
		msg.SetBody("text/plain", *m.Body)

		d := gomail.NewPlainDialer(host, port, from, passwd)
		if skipVerify {
			d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}

		if err := d.DialAndSend(msg); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}
}

func dingFunc() func(*Msg) error {
	corpID := os.Getenv("ESALARM_DING_CORPID")
	secret := os.Getenv("ESALARM_DING_SECRET")
	chatID := os.Getenv("ESALARM_DING_CHATID")

	c := godingtalk.NewDingTalkClient(corpID, secret)
	c.Cache = godingtalk.NewInMemoryCache()
	msgs := make(chan *Msg)
	errs := make(chan error)

	go func() {
		for {
			msg := <-msgs
			if err := c.RefreshAccessToken(); err != nil {
				errs <- errors.WithStack(err)
			}

			if _, err := c.SendTextMessage("", chatID, msg.join("\n\n")); err != nil {
				errs <- errors.WithStack(err)
			}

			errs <- nil
		}
	}()

	return func(m *Msg) error {
		msgs <- m
		return <-errs
	}
}
