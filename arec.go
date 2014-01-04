package main

import (
	"github.com/rakyll/command"
	"github.com/yosisa/arec/reserve"
	"log"
)

func main() {
	command.Parse()

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Load configuration failed: %v", configFile, err)
	}
	log.Printf("Configuration loaded: %s", *configFile)
	reserve.Connect(config.MongoURI)

	command.Run()
}
