package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"bitbucket.org/xcrossing/elastic_alarm/notifiers"
)

type config struct {
	Interval int             `json:"interval"`
	JSON     json.RawMessage `json:"json"`
}

func main() {
	for _, cfg := range configs() {
		// go func(js []byte) {
		// jsonReader := bytes.NewReader(json)
		// fetch(jsonReader)
		jsonReader := bytes.NewReader(cfg.JSON)
		fetch(jsonReader)
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

func fetch(json io.Reader) {
	req, _ := http.NewRequest("GET", "http://192.168.0.107:9200/_search", json)
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
