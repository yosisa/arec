package main

import (
	"flag"
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

	if err := epg.SaveEPG(os.Stdin); err != nil {
		log.Fatal(err)
	}

	rule := reserve.Rule{Keyword: "news"}
	if err := rule.Save(); err != nil {
		log.Print(err)
	}
	reserve.ApplyAllRules(0)
}
