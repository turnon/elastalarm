package notifiers

import (
	"encoding/json"
	"fmt"
)

type stdout struct {
}

func (s *stdout) Send(m *Msg) error {
	fmt.Println(*m.Title, "\n\n", *m.Body)
	return nil
}

func newStdout(cfg json.RawMessage) (Notifier, error) {
	return &stdout{}, nil
}
