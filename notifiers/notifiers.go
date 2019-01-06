package notifiers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hugozhu/godingtalk"
	gomail "gopkg.in/gomail.v2"
)

var Names = make(map[string](func(*Msg)))

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

func dingFunc() func(*Msg) {
	corpID := os.Getenv("ESALARM_DING_CORPID")
	secret := os.Getenv("ESALARM_DING_SECRET")
	chatID := os.Getenv("ESALARM_DING_CHATID")

	c := godingtalk.NewDingTalkClient(corpID, secret)
	c.Cache = godingtalk.NewInMemoryCache()
	msgs := make(chan *Msg)

	go func() {
		for {
			msg := <-msgs
			c.RefreshAccessToken()
			c.SendTextMessage("", chatID, msg.join("\n\n"))
		}
	}()

	return func(m *Msg) {
		msgs <- m
	}
}
