package notifiers

import (
	"crypto/tls"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	gomail "gopkg.in/gomail.v2"
)

type emailSender struct {
	host, passwd, from string
	port               int
	skipVerify         bool
	To                 []string `json:"to"`
}

func (s *emailSender) Send(m *Msg) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.from)
	msg.SetHeader("To", s.To...)
	msg.SetHeader("Subject", *m.Title)
	msg.SetBody("text/plain", *m.Body)

	d := gomail.NewPlainDialer(s.host, s.port, s.from, s.passwd)
	if s.skipVerify {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if err := d.DialAndSend(msg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func newEmail(cfg json.RawMessage) (Notifier, error) {
	sender, err := basicEmail()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(cfg, sender); err != nil {
		return nil, errors.Wrap(err, "收件人配置错误")
	}

	return sender, nil
}

func basicEmail() (*emailSender, error) {
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

	skipVerify := os.Getenv("ESALARM_MAIL_SKIP_VERIFY") != ""

	sender := &emailSender{
		from:       from,
		host:       host,
		passwd:     passwd,
		port:       port,
		skipVerify: skipVerify,
	}

	return sender, nil
}
