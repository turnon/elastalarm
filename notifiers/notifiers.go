package notifiers

import (
	"encoding/json"
	"os"

	"github.com/hugozhu/godingtalk"
	"github.com/pkg/errors"
)

var (
	Names      = make(map[string](func(*Msg) error))
	Errors     = make(map[string]error)
	Generators = map[string](func(cfg json.RawMessage) (Notifier, error)){
		"stdout": newStdout,
		"email":  newEmail,
	}
)

func init() {
	initNotifier("ding", dingFunc)
}

func initNotifier(name string, notifierGenerator func() (func(*Msg) error, error)) {
	notifier, errMsg := notifierGenerator()
	Names[name] = notifier
	Errors[name] = errMsg
}

type Notifier interface {
	Send(*Msg) error
}

type Msg struct {
	Title, Body *string
}

func (msg *Msg) join(seperate string) string {
	return *msg.Title + seperate + *msg.Body
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
