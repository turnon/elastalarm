package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
	"time"

	"github.com/turnon/elastalarm/notifiers"
	"github.com/turnon/elastalarm/paradigms"
)

type config struct {
	Skip         bool            `json:"skip"`
	Title        string          `json:"title"`
	Now          string          `json:"now"`
	TimeF        string          `json:"time_field"`
	Interval     string          `json:"interval"`
	TimeoutRetry int             `json:"timeout_retry"`
	Index        string          `json:"index"`
	ParadigmName string          `json:"paradigm"`
	Condition    json.RawMessage `json:"condition"`
	// Detail       json.RawMessage            `json:"detail"`
	Alarms    map[string]json.RawMessage `json:"alarms"`
	notifiers []notifiers.Notifier
	paradigms.Paradigm
	_reqBody *string
	_ticker  *time.Ticker
}

func loadConfig(path string) *config {
	js, err := ioutil.ReadFile(path)
	if err != nil {
		failToLoad(path, err)
	}

	cfg := &config{}
	if err := json.Unmarshal(js, cfg); err != nil {
		failToLoad(path, err)
	}

	if cfg.Skip {
		return cfg
	}

	cfg.Paradigm = paradigms.Names(cfg.ParadigmName)
	if cfg.Paradigm == nil {
		failToLoad(path, "no such paradigm '"+cfg.ParadigmName+"'")
	}

	if err := json.Unmarshal(cfg.Condition, cfg.Paradigm); err != nil {
		failToLoad(path, err)
	}

	if err := cfg.makeReqBody(); err != nil {
		failToLoad(path, err)
	}

	if err := cfg.makeTicker(); err != nil {
		failToLoad(path, err)
	}

	for notifierName, notifierConfig := range cfg.Alarms {
		gen := notifiers.Generators[notifierName]
		if gen == nil {
			failToLoad(path, "no such generator '"+notifierName+"'")
		}

		notifier, err := gen(notifierConfig)
		if err != nil {
			failToLoad(path, err)
		}

		cfg.notifiers = append(cfg.notifiers, notifier)
	}

	return cfg
}

func failToLoad(path string, errMsg interface{}) {
	panic(fmt.Sprintf("%s : %+v", path, errMsg))
}

func (cfg *config) makeReqBody() error {
	t := template.New("a")
	if _, err := t.Parse(cfg.Template()); err != nil {
		return err
	}

	sb := &strings.Builder{}
	if err := t.Execute(sb, cfg); err != nil {
		return err
	}

	str := sb.String()
	cfg._reqBody = &str

	banner := strings.Repeat("*", len(cfg.Title))
	fmt.Printf("%s\n%s\n%s\n%s\n", banner, cfg.Title, banner, str)

	return nil
}

func (cfg *config) reqBody() *string {
	return cfg._reqBody
}

func (cfg *config) ReqBody() io.Reader {
	return strings.NewReader(*cfg.reqBody())
}

func (cfg *config) NowString() string {
	if cfg.Now == "" || cfg.Now == "now" {
		return "now"
	}
	return cfg.Now + "||"
}

func (cfg *config) TimeField() string {
	if cfg.TimeF == "" || cfg.Now == "@timestamp" {
		return "@timestamp"
	}
	return cfg.TimeF
}

func (cfg *config) makeTicker() error {
	duration, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		return err
	}

	cfg._ticker = time.NewTicker(duration)
	return nil
}

func (cfg *config) stopTicker() {
	cfg._ticker.Stop()
}

func (cfg *config) ticker() <-chan time.Time {
	return cfg._ticker.C
}
