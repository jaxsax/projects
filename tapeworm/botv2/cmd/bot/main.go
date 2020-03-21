package main

import (
	"flag"

	"github.com/jaxsax/projects/tapeworm/botv2"
)

var configPath = flag.String("config_path", "config.yml", "path to config file")

func main() {
	flag.Parse()

	b := botv2.NewBot(*configPath)

	err := b.Init()
	if err != nil {
		panic(err)
	}

	err = b.Run()
	if err != nil {
		panic(err)
	}
}
