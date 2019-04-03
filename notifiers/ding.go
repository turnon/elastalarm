package notifiers

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/hugozhu/godingtalk"
	"github.com/pkg/errors"
)

var (
	dingClientInit    sync.Once
	dingClient        *godingtalk.DingTalkClient
	dingClientInitErr error
	lock              sync.Mutex
)

type ding struct {
	Chats  []string `json:"chats"`
	Users  []string `json:"users"`
	Robots []string `json:"robots"`
}

func (s *ding) Send(m *Msg) error {
	lock.Lock()
	defer lock.Unlock()

	if err := dingClient.RefreshAccessToken(); err != nil {
		return errors.WithStack(err)
	}

	msg := m.join("\n\n")

	for _, chat := range s.Chats {
		if chat == "" {
			continue
		}
		if _, err := dingClient.SendTextMessage("", chat, msg); err != nil {
			return errors.WithStack(err)
		}
	}

	for _, user := range s.Users {
		if user == "" {
			continue
		}
		if err := dingClient.SendAppMessage("", user, msg); err != nil {
			return errors.WithStack(err)
		}
	}

	for _, robot := range s.Robots {
		if robot == "" {
			continue
		}
		if _, err := dingClient.SendRobotTextMessage(robot, msg); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func newDing(cfg json.RawMessage) (Notifier, error) {
	if err := initDingClient(); err != nil {
		return nil, err
	}

	d := &ding{}
	if err := json.Unmarshal(cfg, d); err != nil {
		return nil, errors.Wrap(err, "钉钉配置错误")
	}

	return d, nil

}

func initDingClient() error {
	dingClientInit.Do(func() {
		corpID := os.Getenv("ESALARM_DING_CORPID")
		secret := os.Getenv("ESALARM_DING_SECRET")
		agentID := os.Getenv("ESALARM_DING_AGENT")

		if corpID == "" || secret == "" {
			dingClientInitErr = errors.New("ESALARM_DING_CORPID / ESALARM_DING_SECRET 未设置")
			return
		}

		dingClient = godingtalk.NewDingTalkClient(corpID, secret)
		dingClient.AgentID = agentID
		dingClient.Cache = godingtalk.NewInMemoryCache()

		if err := dingClient.RefreshAccessToken(); err != nil {
			dingClientInitErr = errors.WithStack(err)
			return
		}
	})

	return dingClientInitErr
}
