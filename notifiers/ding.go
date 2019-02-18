package notifiers

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/hugozhu/godingtalk"
	"github.com/pkg/errors"
)

type dingReceivers struct {
	Chats, Users []string
	DingWrapper  *dingWrapper
}

type dingWrapper struct {
	chatMsgs       chan *Msg
	UserMsgs       chan *Msg
	errs           chan error
	dingTalkClient godingtalk.DingTalkClient
}

var (
	parseEnvDingOnce sync.Once
	envDingErr       error
	dingW            dingWrapper
)

func initDing(config *json.RawMessage) (notifier, error) {
	if err := initCommonEmail(); err != nil {
		return nil, envEmailErr
	}

	d := &dingReceivers{}
	if err := json.Unmarshal(*config, d); err != nil {
		return nil, errors.Wrap(err, "钉钉接收者配置有误")
	}
	d.DingWrapper = &dingW

	return d, nil
}

func initCommonDing() error {
	parseEnvDingOnce.Do(func() {
		corpID := os.Getenv("ESALARM_DING_CORPID")
		secret := os.Getenv("ESALARM_DING_SECRET")

		if corpID == "" || secret == "" {
			envDingErr = errors.New("ESALARM_DING_CORPID / ESALARM_DING_SECRET 未设置")
			return
		}

		c := godingtalk.NewDingTalkClient(corpID, secret)
		c.Cache = godingtalk.NewInMemoryCache()

		if err := c.RefreshAccessToken(); err != nil {
			envDingErr = errors.WithStack(err)
			return
		}

		chatMsgs := make(chan *Msg)
		userMsgs := make(chan *Msg)
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

		dingW = dingWrapper{chatMsgs, userMsgs, errs, *c}
	})

	return envDingErr
}

func (d *dingReceivers) Send(m *Msg) error {
	for _, chat := range d.Chats {
		d.sendChat(chat, m)
	}
	return nil
}

func (d *dingReceivers) sendChat(chatID string, msg *Msg) error {
	if _, err := d.DingWrapper.dingTalkClient.SendTextMessage("", chatID, msg.join("\n\n")); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func dingFunc() (func(*Msg) error, error) {

	chatID := os.Getenv("ESALARM_DING_CHATID")
	msgs := make(chan *Msg)
	errs := make(chan error)

	return func(m *Msg) error {
		msgs <- m
		return <-errs
	}, nil
}
