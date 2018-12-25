package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"bitbucket.org/xcrossing/elastic_alarm/notifiers"
)

type config struct {
	Interval int             `json:"interval"`
	JSON     json.RawMessage `json:"json"`
}

func main() {

	host := flag.String("host", "http://0.0.0.0:9200", "es host")
	flag.Parse()

	for _, cfg := range configs() {
		// go func(js []byte) {
		// jsonReader := bytes.NewReader(json)
		// fetch(jsonReader)
		jsonReader := bytes.NewReader(cfg.JSON)
		fetch(jsonReader, *host)
		// }(js)
	}
	// <-make(chan bool)
}

func configs() []config {
	var configArray []config

	files, err := filepath.Glob("configs/*")
	if err != nil {
		panic(err)
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
		configArray = append(configArray, cfg)
	}

	return configArray
}

func fetch(json io.Reader, host string) {
	var sb strings.Builder
	sb.WriteString(host)
	sb.WriteString("/_search")
	url := sb.String()

	req, _ := http.NewRequest("GET", url, json)
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

	n := notifiers.Email{}
	n.SetTitle(resp.Status)
	n.SetBody(string(body))
	notify(&n)

}

func notify(n notifiers.Notifier) {
	n.Notify()
}
