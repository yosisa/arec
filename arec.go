package main

import (
	"github.com/voxelbrain/goptions"
	"github.com/yosisa/arec/reserve"
	"log"
	"os"
)

func main() {
	options := CmdOptions{
		Config: "./arec.json",
	}
	goptions.ParseAndFail(&options)

	command, ok := commands[string(options.Verbs)]
	if !ok {
		goptions.PrintHelp()
		os.Exit(1)
	}

	config, err := LoadConfig(&options.Config)
	reserve.Recpt1Path = config.Recpt1
	if err != nil {
		log.Fatalf("Load configuration from %s failed: %v", options.Config, err)
	}
	log.Printf("Configuration loaded from %s", options.Config)
	reserve.Connect(config.MongoURI)

	command(&options, &config)
}
