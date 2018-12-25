package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/xcrossing/elastic_alarm/notifiers"
)

var (
	intervalRe = regexp.MustCompile("(\\d+)(.)")
)

type config struct {
	Title    string          `json:"title"`
	Interval string          `json:"interval"`
	JSON     json.RawMessage `json:"json"`
}

func (cfg *config) load(path string) {
	js, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(js, cfg); err != nil {
		panic(err)
	}
}

func (cfg *config) ticker() <-chan time.Time {
	duration, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		panic(err)
	}
	return time.NewTicker(duration).C
}

func (cfg *config) monitor(host string) {
	go func() {
		for range cfg.ticker() {
			cfg.fetch(host)
		}
	}()
}

func (cfg *config) requestBody() io.Reader {
	return bytes.NewReader(cfg.JSON)
}

func (cfg *config) fetch(host string) {
	var sb strings.Builder
	sb.WriteString(host)
	sb.WriteString("/_search")
	url := sb.String()

	req, _ := http.NewRequest("GET", url, cfg.requestBody())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	n := notifiers.Stdout{}
	n.SetTitle(cfg.Title)
	n.SetBody(string(body))
	n.Notify()

}
