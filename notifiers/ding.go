package notifiers

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/hugozhu/godingtalk"
	"github.com/pkg/errors"
)

var (
	dingClientInit    sync.Once
	dingClient        godingtalk.DingTalkClient
	dingClientInitErr error
)

type ding struct {
}

func (s *ding) Send(m *Msg) error {
	fmt.Println(*m.Title, "\n\n", *m.Body)
	return nil
}

func newDing(cfg json.RawMessage) (Notifier, error) {
	if err := initDingClient(); err != nil {
		return nil, err
	}

	return nil, dingClientInitErr

}

func initDingClient() error {
	dingClientInit.Do(func() {
		corpID := os.Getenv("ESALARM_DING_CORPID")
		secret := os.Getenv("ESALARM_DING_SECRET")

		if corpID == "" || secret == "" {
			dingClientInitErr = errors.New("ESALARM_DING_CORPID / ESALARM_DING_SECRET 未设置")
			return
		}

		dingClient := godingtalk.NewDingTalkClient(corpID, secret)
		dingClient.Cache = godingtalk.NewInMemoryCache()

		if err := dingClient.RefreshAccessToken(); err != nil {
			dingClientInitErr = errors.WithStack(err)
			return
		}
	})

	return dingClientInitErr
}
