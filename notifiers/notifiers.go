package notifiers

import (
	"encoding/json"
	"fmt"

	"github.com/turnon/elastalarm/response"
)

var (
	Generators = map[string](func(cfg json.RawMessage) (Notifier, error)){
		"stdout": newStdout,
		"email":  newEmail,
		"ding":   newDing,
	}
)

type Notifier interface {
	Send(*Msg) error
}

type Msg struct {
	Title string
	*response.Result
}

func (msg *Msg) TextWithTitle() string {
	return fmt.Sprintf("%s\n\n%s", msg.Title, msg.Text())
}
