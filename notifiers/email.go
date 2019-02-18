package notifiers

import (
	"crypto/tls"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	gomail "gopkg.in/gomail.v2"
)

type envMail struct {
	host, passwd, from string
	port               int
	skipVerify         bool
}

type email struct {
	To []string
}

type emailReceivers struct {
	To []string
}

var (
	parseEnvMailOnce sync.Once
	envEmailErr      error
	envEmailCfg      envMail
)

func initEmail(config *json.RawMessage) (notifier, error) {
	if err := initCommonEmail(); err != nil {
		return nil, envEmailErr
	}

	e := &email{}
	if err := json.Unmarshal(*config, e); err != nil {
		return nil, errors.Wrap(err, "邮件接收者配置有误")
	}

	return e, nil
}

func initCommonEmail() error {
	parseEnvMailOnce.Do(func() {
		common := envMail{}
		server := os.Getenv("ESALARM_MAIL_SERVER")
		common.from = os.Getenv("ESALARM_MAIL_FROM")
		common.passwd = os.Getenv("ESALARM_MAIL_PASSWD")
		if server == "" || common.from == "" || common.passwd == "" {
			envEmailErr = errors.New("邮件服务未配置 ESALARM_MAIL_SERVER / ESALARM_MAIL_FROM / ESALARM_MAIL_PASSWD")
			return
		}

		hostPort := strings.Split(server, ":")
		common.host = hostPort[0]
		port, err := strconv.Atoi(hostPort[1])
		if err != nil {
			envEmailErr = errors.Wrap(err, "无法解析 ESALARM_MAIL_SERVER")
			return
		}
		common.port = port

		common.skipVerify = os.Getenv("ESALARM_MAIL_SKIP_VERIFY") != ""

		envEmailCfg = common
	})

	return envEmailErr
}

func (e *email) Send(m *Msg) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", envEmailCfg.from)
	msg.SetHeader("To", e.To...)
	msg.SetHeader("Subject", *m.Title)
	msg.SetBody("text/plain", *m.Body)

	d := gomail.NewPlainDialer(envEmailCfg.host, envEmailCfg.port, envEmailCfg.from, envEmailCfg.passwd)
	if envEmailCfg.skipVerify {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if err := d.DialAndSend(msg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
