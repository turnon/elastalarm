package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	host := os.Getenv("ESALARM_HOST")
	if host == "" {
		panic("ES地址未配置")
	}

	cfgPath := flag.String("configs", "config", "config files location")
	flag.Parse()

	initMonitors(host, configFiles(*cfgPath))

	<-make(chan bool)
}

func configFiles(dir string) []string {
	files, err := filepath.Glob(dir + "/*")

	if err != nil {
		panic(err)
	}

	if len(files) == 0 {
		panic("config not found")
	}

	return files
}
