package notifiers

import (
	"encoding/json"
	"fmt"

	"github.com/turnon/elastalarm/response"
)

var (
	Generators = map[string](func(cfg json.RawMessage) (Notifier, error)){
		"stdout":   newStdout,
		"email":    newEmail,
		"ding":     newDing,
		"web_hook": newWebHook,
	}
)

type Notifier interface {
	Send(*Msg) error
}

type Msg struct {
	Title string `json:"title"`
	*response.Result
}

func (msg *Msg) TextWithTitle() string {
	return fmt.Sprintf("%s\n\n%s", msg.Title, msg.Text())
}

func (msg *Msg) JSON() ([]byte, error) {
	return json.Marshal(msg)
}
