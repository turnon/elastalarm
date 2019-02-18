package notifiers

import (
	"encoding/json"
)

var (
	Names = make(map[string](func(config *json.RawMessage) (notifier, error)))
)

func init() {
	Names["stdout"] = initStdout
	Names["email"] = initEmail
	Names["ding"] = initDing
}

// func initNotifier(name string, notifierGenerator func() (func(*Msg) error, error)) {
// 	notifier, errMsg := notifierGenerator()
// 	Names[name] = notifier
// 	Errors[name] = errMsg
// }

type notifier interface {
	Send(m *Msg) error
}

type Msg struct {
	Title, Body *string
}

func (msg *Msg) join(seperate string) string {
	return *msg.Title + seperate + *msg.Body
}
