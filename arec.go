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
	if err != nil {
		log.Fatalf("Load configuration failed: %v", options.Config, err)
	}
	log.Printf("Configuration loaded: %s", &options.Config)
	reserve.Connect(config.MongoURI)

	command(&options)
}
