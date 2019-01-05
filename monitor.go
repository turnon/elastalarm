package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"bitbucket.org/xcrossing/elastic_alarm/notifiers"
	"bitbucket.org/xcrossing/elastic_alarm/response"
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
	url := strings.Join([]string{host, cfg.Index, "_search"}, "/")
	return &monitor{config: cfg, url: url, httpClient: &http.Client{}}
}

func (mon *monitor) run() {
	go func() {
		mon.check()
		for range mon.ticker() {
			mon.check()
		}
	}()
}

func (mon *monitor) check() {
	req, _ := http.NewRequest("GET", mon.url, mon.ReqBody())
	req.Header.Set("Content-Type", "application/json")

	resp, err := mon.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	respObj := &response.Response{}
	respObj.Unmarshal(body)

	found, detail := mon.Found(respObj)
	if !found {
		return
	}

	msg := notifiers.Msg{Title: &mon.Title, Body: detail}

	for notifier, targets := range mon.Alarms {
		notifiers.Names[notifier](&msg, &targets)
	}
}
