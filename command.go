package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/yosisa/arec/command"
	"github.com/yosisa/arec/epg"
	"github.com/yosisa/arec/reserve"
)

const (
	GR_REC_TIME = 90 * time.Second
	BS_REC_TIME = 300 * time.Second
)

var app *cli.App

func init() {
	app = cli.NewApp()
	app.Name = "arec"
	app.Usage = "Japanese TV recorder"
	app.Flags = []cli.Flag{
		cli.StringFlag{"config, c", "arec.json", "Config file"},
	}
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Run scheduler",
			Action: CmdScheduler,
		},
		{
			Name:  "rule",
			Usage: "Manage auto reservation rules",
			Flags: []cli.Flag{
				cli.StringFlag{"title, t", "", "Regexp for title"},
			},
			Action: CmdRule,
		},
		{
			Name:  "epg",
			Usage: "Update EPG data",
			Flags: []cli.Flag{
				cli.StringFlag{"ch", "", "Update a given channel only"},
				cli.StringFlag{"json", "", "Feed from a given json file"},
			},
			Action: CmdEPG,
		},
	}
}

func CmdScheduler(c *cli.Context) {
	engine := reserve.NewEngine(2, 0)
	engine.ReserveFromDB()
	engine.RunForever(engine.ReserveFromDB)
}

func CmdRule(c *cli.Context) {
	title := c.String("title")
	if title == "" {
		log.Fatal("One or more conditions needed.")
	}
	rule := reserve.Rule{Keyword: title}
	if err := rule.Save(); err != nil {
		log.Fatal(err)
	}
	rule.Apply(0)
}

func CmdEPG(c *cli.Context) {
	if path := c.String("json"); path != "" {
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		if err := epg.SaveEPG(f, c.String("ch")); err != nil {
			log.Fatal(err)
		}
	}

	config := LoadConfigAndInit(c)
	engine := reserve.NewEngine(2, 2)
	updateEPG := func() {
		for _, channel := range config.Channels["GR"] {
			epg.Reserve(engine, "GR", channel.Ch)
		}

		if bs, ok := config.Channels["BS"]; ok && len(bs) > 0 {
			epg.Reserve(engine, "BS", bs[0].Ch)
		}
	}
	updateEPG()

	engine.RunForever(updateEPG)
}

func LoadConfigAndInit(c *cli.Context) *Config {
	path := c.GlobalString("config")
	config, err := LoadConfig(path)
	if err != nil {
		log.Fatalf("Load configuration %s failed: %v", path, err)
	}

	log.Printf("Configuration loaded from %s", path)
	command.Recpt1Path = config.Recpt1
	command.EpgdumpPath = config.Epgdump
	reserve.Connect(config.MongoURI)

	return config
}
