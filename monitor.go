package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/turnon/elastalarm/notifiers"
	"github.com/turnon/elastalarm/response"
)

type monitor struct {
	httpClient *http.Client
	url        string
	*config
}

func initMonitors(host string, files []string) {
	for _, file := range files {
		cfg, err := loadConfig(file)
		if err != nil {
			log.Println(err)
			continue
		}

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
	if resp.StatusCode != 200 {
		log.Printf("%d|%s|%s\n", resp.StatusCode, mon.Title, string(body))
		return
	}

	respObj := &response.Response{}
	respObj.Unmarshal(body)

	found, detail := mon.Found(respObj)
	if !found {
		return
	}

	msg := notifiers.Msg{Title: &mon.Title, Body: detail}

	for _, notifier := range mon.Alarms {
		notifiers.Names[notifier](&msg)
	}
}
