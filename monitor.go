package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/turnon/elastalarm/notifiers"
	"github.com/turnon/elastalarm/response"
)

type monitor struct {
	httpClient *http.Client
	url        string
	done       chan bool
	*config
}

const timeOut = 3 * time.Second

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
	return &monitor{
		config:     cfg,
		url:        url,
		httpClient: &http.Client{Timeout: timeOut},
		done:       make(chan bool),
	}
}

func (mon *monitor) run() {
	go func() {
		mon.check()

		for {
			select {
			case <-mon.ticker():
				mon.check()
			case <-mon.done:
				return
			}
		}
	}()
}

func (mon *monitor) check() {
	err := mon._check()
	if err != nil {
		mon.stopTicker()
		close(mon.done)
		log.Printf("%+v", err)
	}
}

func (mon *monitor) _check() error {
	req, err := http.NewRequest("GET", mon.url, mon.ReqBody())
	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := mon.httpClient.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status + "|" + string(body))
	}

	respObj := &response.Response{}
	if err := respObj.Unmarshal(body); err != nil {
		return errors.WithStack(err)
	}

	found, detail := mon.Found(respObj)
	if !found {
		return nil
	}

	msg := notifiers.Msg{Title: &mon.Title, Body: detail}

	for _, notifier := range mon.Alarms {
		if err := notifiers.Names[notifier](&msg); err != nil {
			log.Printf("%+v", err)
		}
	}

	return nil
}
