package notifiers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hugozhu/godingtalk"
	"github.com/pkg/errors"
)

var (
	Names  = make(map[string](func(*Msg) error))
	Errors = make(map[string]error)
)

func init() {
	Names["stdout"] = stdout
	initNotifier("email", emailFunc)
	initNotifier("ding", dingFunc)
}

func initNotifier(name string, notifierGenerator func() (func(*Msg) error, error)) {
	notifier, errMsg := notifierGenerator()
	Names[name] = notifier
	Errors[name] = errMsg
}

type Msg struct {
	Title, Body *string
}

func (msg *Msg) join(seperate string) string {
	return *msg.Title + seperate + *msg.Body
}

type notifier interface {
	Send(*Msg) error
}

func stdout(m *Msg) error {
	fmt.Println(*m.Title, "\n\n", *m.Body)
	return nil
}

func emailFunc() (map[string]string, error) {
	server := os.Getenv("ESALARM_MAIL_SERVER")
	from := os.Getenv("ESALARM_MAIL_FROM")
	passwd := os.Getenv("ESALARM_MAIL_PASSWD")
	if server == "" || from == "" || passwd == "" {
		return nil, errors.New("邮件服务未配置 ESALARM_MAIL_SERVER / ESALARM_MAIL_FROM / ESALARM_MAIL_PASSWD")
	}

	hostPort := strings.Split(server, ":")
	host := hostPort[0]
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return nil, errors.Wrap(err, "无法解析 ESALARM_MAIL_SERVER")
	}
	// to := os.Getenv("ESALARM_MAIL_TO")
	to := []string{"yuanzp@reocar.com", "1020715764@qq.com"}
	skipVerify := os.Getenv("ESALARM_MAIL_SKIP_VERIFY") != ""

	mail := &email{}

	// emailFunction := func(m *Msg) error {
	// 	msg := gomail.NewMessage()
	// 	msg.SetHeader("From", from)
	// 	msg.SetHeader("To", to...)
	// 	msg.SetHeader("Subject", *m.Title)
	// 	msg.SetBody("text/plain", *m.Body)

	// 	d := gomail.NewPlainDialer(host, port, from, passwd)
	// 	if skipVerify {
	// 		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// 	}

	// 	if err := d.DialAndSend(msg); err != nil {
	// 		return errors.WithStack(err)
	// 	}

	// 	return nil
	// }

	// return emailFunction, nil
}

func dingFunc() (func(*Msg) error, error) {
	corpID := os.Getenv("ESALARM_DING_CORPID")
	secret := os.Getenv("ESALARM_DING_SECRET")

	if corpID == "" || secret == "" {
		return nil, errors.New("ESALARM_DING_CORPID / ESALARM_DING_SECRET 未设置")
	}

	c := godingtalk.NewDingTalkClient(corpID, secret)
	c.Cache = godingtalk.NewInMemoryCache()

	if err := c.RefreshAccessToken(); err != nil {
		return nil, errors.WithStack(err)
	}

	chatID := os.Getenv("ESALARM_DING_CHATID")
	msgs := make(chan *Msg)
	errs := make(chan error)

	go func() {
		for {
			msg := <-msgs
			if err := c.RefreshAccessToken(); err != nil {
				errs <- errors.WithStack(err)
				continue
			}

			if _, err := c.SendTextMessage("", chatID, msg.join("\n\n")); err != nil {
				errs <- errors.WithStack(err)
				continue
			}

			errs <- nil
		}
	}()

	return func(m *Msg) error {
		msgs <- m
		return <-errs
	}, nil
}
