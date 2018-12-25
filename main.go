package main

import (
	"flag"
	"path/filepath"
)

func main() {
	host := flag.String("host", "http://0.0.0.0:9200", "es host")
	cfgPath := flag.String("configs", "config", "config files location")
	flag.Parse()

	for _, cfg := range configs(*cfgPath) {
		(func(cfg config) {
			cfg.monitor(*host)
		})(cfg)

	}

	<-make(chan bool)
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
		cfg := config{}
		cfg.load(file)
		configArray = append(configArray, cfg)
	}

	return configArray
}
