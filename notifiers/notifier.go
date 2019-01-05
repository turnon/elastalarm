package notifiers

import (
	"fmt"
)

// Notifier interface
type Notifier interface {
	SetTitle(s string)
	SetBody(s string)
	Notify()
}

type Msg struct {
	Title, Body *string
}

var Names = make(map[string](func(*Msg, *[]string)))

func init() {
	Names["stdout"] = stdout
}

func stdout(m *Msg, targets *[]string) {
	fmt.Println(*m.Title, *m.Body)
}
