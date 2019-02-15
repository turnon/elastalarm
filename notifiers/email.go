package notifiers

type email struct {
	host, from, passwd string
	port               int
	To                 string `json:"to"`
}

func (m *email) Send(payload *Msg) error {

	return nil
}
