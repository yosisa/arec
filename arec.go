package main

import (
	"flag"
	"fmt"
	"github.com/yosisa/arec/epg"
	"github.com/yosisa/arec/reserve"
	"log"
	"os"
)

var configFile *string = flag.String("config", "./arec.json", "path to config file")

func main() {
	flag.Parse()

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Load configuration failed: %v", configFile, err)
	}
	log.Printf("Configuration loaded: %s", *configFile)
	reserve.Connect(config.MongoURI)

	data, err := epg.DecodeJson(os.Stdin)
	if err != nil {
		panic(err)
	}
	channel := data[0]
	if err := channel.Save(); err != nil {
		fmt.Printf("%+v\n", err)
	}

	for _, program := range channel.Programs {
		if err := program.Save(); err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}
