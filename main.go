package main

import (
	"flag"
	"path/filepath"
)

func main() {
	host := flag.String("host", "http://0.0.0.0:9200", "es host")
	cfgPath := flag.String("configs", "config", "config files location")
	flag.Parse()

	initMonitors(*host, configFiles(*cfgPath))

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

// func main() {
// 	cfg := loadConfig(".\\configs\\aggs.1.json")
// 	fmt.Println(*cfg.Percentage.ReqBody())
// 	// fmt.Println(paradigms.Percentage{} == s)
// }
