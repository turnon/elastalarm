package notifiers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type webHook struct {
	*http.Client
	Url    string `json:"url"`
	Method string `json:"method"`
}

const timeOut = 3 * time.Second

func (wh *webHook) Send(m *Msg) error {
	data, err := m.JSON()
	if err != nil {
		return errors.WithStack(err)
	}

	body := bytes.NewReader(data)
	req, err := http.NewRequest(wh.Method, wh.Url, body)
	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := wh.Client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	return nil
}

func newWebHook(cfg json.RawMessage) (Notifier, error) {
	wh := &webHook{}
	if err := json.Unmarshal(cfg, wh); err != nil {
		return nil, errors.Wrap(err, "web hook配置错误")
	}

	if wh.Method == "" {
		wh.Method = "POST"
	}

	wh.Client = &http.Client{Timeout: timeOut}
	return wh, nil
}
