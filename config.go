package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
	"time"
)

var (
	intervalRe = regexp.MustCompile("(\\d+)(.)")
)

type config struct {
	Title    string          `json:"title"`
	Interval string          `json:"interval"`
	JSON     json.RawMessage `json:"json"`
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
	return bytes.NewReader(cfg.JSON)
}
