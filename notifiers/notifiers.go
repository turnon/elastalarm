package notifiers

import (
	"encoding/json"
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
	Title, Body *string
}

func (msg *Msg) join(seperate string) string {
	return *msg.Title + seperate + *msg.Body
}
