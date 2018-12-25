package notifiers

import (
	"fmt"
)

// Stdout to print
type Stdout struct {
	title string
	body  string
}

// SetTitle set title
func (o *Stdout) SetTitle(s string) {
	o.title = s
}

// SetBody set body
func (o *Stdout) SetBody(s string) {
	o.body = s
}

// Notify print
func (o *Stdout) Notify() {
	fmt.Println(o.title, o.body)
}
