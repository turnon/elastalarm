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
	cfgPath := flag.String("configs", "config", "config files location")
	flag.Parse()

	for _, cfg := range configs(*cfgPath) {
		// go func(js []byte) {
		// jsonReader := bytes.NewReader(json)
		// fetch(jsonReader)
		jsonReader := bytes.NewReader(cfg.JSON)
		fetch(jsonReader, *host)
		// }(js)
	}
	// <-make(chan bool)
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

	n := notifiers.Stdout{}
	n.SetTitle(resp.Status)
	n.SetBody(string(body))
	notify(&n)

}

func notify(n notifiers.Notifier) {
	n.Notify()
}
