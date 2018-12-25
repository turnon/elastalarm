package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/xcrossing/elastic_alarm/notifiers"
)

var (
	intervalRe = regexp.MustCompile("(\\d+)(.)")
)

type config struct {
	Interval string          `json:"interval"`
	JSON     json.RawMessage `json:"json"`
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

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	n := notifiers.Stdout{}
	n.SetTitle(resp.Status)
	n.SetBody(string(body))
	n.Notify()

}

func configs(path string) []config {
	var configArray []config
	files, err := filepath.Glob(path + "/*")

	if err != nil {
		panic(err)
	}

	if len(files) == 0 {
		panic("config not found")
	}

	for _, file := range files {
		js, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		cfg := config{}
		if err = json.Unmarshal(js, &cfg); err != nil {
			panic(err)
		}
		cfg.ticker()
		configArray = append(configArray, cfg)
	}

	return configArray
}
