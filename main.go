package main

import (
	"flag"
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
