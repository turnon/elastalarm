package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"bitbucket.org/xcrossing/elastic_alarm/paradigms"
)

type config struct {
	Skip      bool            `json:"skip"`
	Title     string          `json:"title"`
	Interval  string          `json:"interval"`
	Index     string          `json:"index"`
	Paradigm  string          `json:"paradigm"`
	Condition json.RawMessage `json:"condition"`
}

func loadConfig(path string) *config {
	js, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	cfg := &config{}
	if err = json.Unmarshal(js, cfg); err != nil {
		panic(err)
	}

	return cfg
}

func (cfg *config) ticker() <-chan time.Time {
	duration, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		panic(err)
	}
	return time.NewTicker(duration).C
}

func (cfg *config) requestBody() io.Reader {
	var paradigm paradigms.Paradigm
	if cfg.Paradigm == "percentage" {
		paradigm = &paradigms.Percentage{}
	}

	if err := json.Unmarshal(cfg.Condition, paradigm); err != nil {
		panic(err)
	}

	return paradigm.ReqBody()
	// return bytes.NewReader(cfg.JSON)
}
