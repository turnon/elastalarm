package notifiers_test

import (
	"os"
	"testing"
	"time"

	"github.com/turnon/elastalarm/notifiers"
)

var (
	t     = time.Now().String()
	title = "title:" + t
	body  = "body:" + t
	msg   = notifiers.Msg{Title: &title, Body: &body}
)

func dingConfig() []byte {
	chat := os.Getenv("ESALARM_DING_TEST_CHAT")
	user := os.Getenv("ESALARM_DING_TEST_USER")
	robot := os.Getenv("ESALARM_DING_TEST_ROBOT")
	return []byte("{\"chats\": [\"" + chat + "\"], \"users\": [\"" + user + "\"], \"robots\": [\"" + robot + "\"]}")
}

func TestDing(t *testing.T) {
	dingGenerator := notifiers.Generators["ding"]

	ding, err := dingGenerator(dingConfig())
	if err != nil {
		t.Error(err)
	}

	if err := ding.Send(&msg); err != nil {
		t.Error(err)
	}
}

func emailConfig() []byte {
	to := os.Getenv("ESALARM_MAIL_TEST_TO")
	return []byte("{\"to\": [\"" + to + "\"]}")
}

func TestEmail(t *testing.T) {
	emailGenerator := notifiers.Generators["email"]

	email, err := emailGenerator(emailConfig())
	if err != nil {
		t.Error(err)
	}

	if err := email.Send(&msg); err != nil {
		t.Error(err)
	}
}
