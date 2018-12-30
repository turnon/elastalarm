package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"bitbucket.org/xcrossing/elastic_alarm/notifiers"
)

type monitor struct {
	httpClient *http.Client
	url        string
	*config
}

func initMonitors(host string, files []string) {
	for _, file := range files {
		cfg := loadConfig(file)
		if cfg.Skip {
			continue
		}
		m := newMonitor(host, cfg)
		m.run()
	}
}

func newMonitor(host string, cfg *config) *monitor {
	var sb strings.Builder
	sb.WriteString(host)
	sb.WriteString("/_search")
	url := sb.String()

	return &monitor{config: cfg, url: url, httpClient: &http.Client{}}
}

func (m *monitor) run() {
	go func() {
		for range m.ticker() {
			m.check()
		}
	}()
}

func (m *monitor) check() {
	req, _ := http.NewRequest("GET", m.url, m.requestBody())
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	n := notifiers.Stdout{}
	n.SetTitle(m.Title)
	n.SetBody(string(body))
	n.Notify()

}
