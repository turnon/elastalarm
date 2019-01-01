package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
	"time"

	"bitbucket.org/xcrossing/elastic_alarm/paradigms"
)

type config struct {
	Skip         bool            `json:"skip"`
	Title        string          `json:"title"`
	Interval     string          `json:"interval"`
	Index        string          `json:"index"`
	ParadigmName string          `json:"paradigm"`
	Condition    json.RawMessage `json:"condition"`
	paradigms.Paradigm
	_reqBody *string
}

func loadConfig(path string) *config {
	js, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	cfg := &config{}
	if err := json.Unmarshal(js, cfg); err != nil {
		panic(err)
	}

	if cfg.Skip {
		return cfg
	}

	if cfg.ParadigmName == "percentage" {
		cfg.Paradigm = &paradigms.Percentage{}
	}

	if err := json.Unmarshal(cfg.Condition, cfg.Paradigm); err != nil {
		panic(err)
	}

	return cfg
}

func (cfg *config) reqBody() *string {
	if cfg._reqBody == nil {
		t := template.New("a")
		t.Parse(cfg.Template())
		sb := &strings.Builder{}
		t.Execute(sb, cfg)
		str := sb.String()
		cfg._reqBody = &str
	}

	return cfg._reqBody
}

func (cfg *config) ReqBody() io.Reader {
	return strings.NewReader(*cfg.reqBody())
}

func (cfg *config) ticker() <-chan time.Time {
	duration, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		panic(err)
	}
	return time.NewTicker(duration).C
}
