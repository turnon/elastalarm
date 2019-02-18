package notifiers

import (
	"encoding/json"
	"fmt"
)

type stdout struct {
}

func initStdout(config *json.RawMessage) (notifier, error) {
	return &stdout{}, nil
}

func (s *stdout) Send(m *Msg) error {
	fmt.Println(*m.Title, "\n\n", *m.Body)
	return nil
}
